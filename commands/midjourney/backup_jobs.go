package midjourney

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jimeh/go-midjourney"
	"github.com/jimeh/mje/commands/shared"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewBackupJobs(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "backup-jobs",
		Aliases: []string{"backup"},
		Short:   "Backup recent jobs data and images",
		RunE:    midjourneyBackupJobsRunE(mc),
	}

	cmd.Flags().Int("limit", -1, "limit of jobs to list")

	cmd.Flags().StringP("type", "t", "", "type of jobs to list")
	cmd.Flags().StringP("order", "o", "new", "either \"new\" or \"oldest\"")
	cmd.Flags().StringP("user-id", "u", "", "user ID to list jobs for")
	cmd.Flags().Bool("fetch-all-jobs", true, "fetch all jobs")
	cmd.Flags().IntP("page", "p", 0, "page to fetch")
	cmd.Flags().StringP("prompt", "s", "", "prompt text to search for")
	cmd.Flags().Int("max-concurrent-downloads", 10, "max concurrent downloads")
	cmd.Flags().Bool("dedupe", true, "dedupe results")
	cmd.Flags().Bool("cached-jobs", false, "use cached all-jobs.json")
	cmd.Flags().String("output", "backup-jobs", "output directory defaults to ./backup-jobs")

	return cmd, nil
}

const MaxWorkers = 5

func midjourneyBackupJobsRunE(mc *midjourney.Client) shared.RunEFunc {
	return func(cmd *cobra.Command, _ []string) error {
		fs := cmd.Flags()
		q := &midjourney.RecentJobsQuery{}

		if v, err := fs.GetInt("amount"); err == nil && v > 0 {
			q.Amount = v
		} else {
			q.Amount = 50
		}

		limit := -1
		if v, err := fs.GetInt("limit"); err == nil && v > 0 {
			limit = v
		}

		if v, err := fs.GetString("type"); err == nil && v != "" {
			q.JobType = midjourney.JobType(v)
		}
		if v, err := fs.GetString("order"); err == nil && v != "" {
			q.OrderBy = midjourney.Order(v)
		}
		if v, err := fs.GetString("user-id"); err == nil && v != "" {
			q.UserID = v
		} else if limit <= 0 {
			return fmt.Errorf("cannot specify unlimited --limit and no --user-id specified")
		}

		if v, err := fs.GetInt("page"); err == nil && v != 0 {
			q.Page = v
		}
		if v, err := fs.GetString("prompt"); err == nil && v != "" {
			q.Prompt = v
		}
		if v, err := fs.GetBool("dedupe"); err == nil {
			q.Dedupe = v
		}

		maxDownloadWorkers := 10
		if v, err := fs.GetInt("max-concurrent-downloads"); err == nil && v != 0 {
			maxDownloadWorkers = v
		}

		fetchAllJobs := true
		var rj *midjourney.RecentJobs

		if v, err := fs.GetBool("fetch-all-jobs"); err == nil {
			if fetchAllJobs && q.UserID == "" {
				return fmt.Errorf("fetch-all-jobs only valid when fetching recent jobs of --user-id")
			}
			fetchAllJobs = v
		}

		rjs := []*midjourney.RecentJobs{}
		r := []*midjourney.Job{}
		totalJobs := 0

		output := shared.FlagString(cmd, "output")
		// Create backup dir
		outputDir := output
		if q.UserID != "" {
			outputDir = filepath.Join(outputDir, q.UserID)
		} else if q.Prompt != "" {
			outputDir = filepath.Join(outputDir, q.Prompt)
		} else {
			outputDir = filepath.Join(outputDir, "recent")
		}

		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return err
		}

		jobsJSONFile := filepath.Join(outputDir, "all-jobs.json")
		if v, err := fs.GetBool("cached-jobs"); err == nil && v {
			allJobsJSON, err := os.ReadFile(jobsJSONFile)
			if err != nil {
				return fmt.Errorf("failed to read all-jobs.json file %s: %w", jobsJSONFile, err)
			}
			err = json.Unmarshal(allJobsJSON, &r)
			if err != nil {
				return fmt.Errorf("failed to unmarshal all-jobs.json file %s : %w", jobsJSONFile, err)
			}
		} else {
			for hasMoreJobs := true; hasMoreJobs; hasMoreJobs = fetchAllJobs && (rj == nil || len(rj.Jobs) >= q.Amount) {
				var err error
				rj, err = mc.RecentJobs(cmd.Context(), q)
				if err != nil {
					return err
				}
				totalJobs += len(rj.Jobs)
				log.Infof("Fetched recent jobs user-id=%s page=%d total-jobs=%d", q.UserID, q.Page, totalJobs)
				q.Page++

				if limit > 0 && totalJobs >= limit {
					break
				}

				rjs = append(rjs, rj)
			}

			for _, rj = range rjs {
				r = append(r, rj.Jobs...)
			}
		}

		imgsDir := filepath.Join(outputDir, "images")
		err = os.MkdirAll(imgsDir, os.ModePerm)
		if err != nil {
			return err
		}

		jobsJSON, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal jobs JSON %w", err)
		}
		// Write out jobs JSON
		os.WriteFile(jobsJSONFile, jobsJSON, 0644)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Hour)
		defer cancel()

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(maxDownloadWorkers)

		for _, j := range r {
			job := j
			g.Go(func() error {
				// Check if directory exists
				jobDir := filepath.Join(outputDir, job.ID)
				err := os.MkdirAll(jobDir, os.ModePerm)
				if err != nil {
					return err
				}

				jobJSONFile := filepath.Join(jobDir, "job.json")
				if !fileExists(jobJSONFile) {
					jobJSON, err := json.MarshalIndent(job, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal job JSON jobID=%s %w", job.ID, err)
					}
					// Write out job JSON
					err = os.WriteFile(jobJSONFile, jobJSON, 0644)
					if err != nil {
						return fmt.Errorf("failed to write job file =%s %w", jobJSONFile, err)
					}
				}

				// Download filepaths if doesn't exist
				for _, img := range job.ImagePaths {
					imgUrl := img // https://golang.org/doc/faq#closures_and_goroutines

					imgFileName := path.Base(imgUrl)

					// Save just image
					createdAt := job.EnqueueTime.Local()
					imgsImgFileOutput := filepath.Join(imgsDir, fmt.Sprintf("%s_%s_%s", createdAt.Format(time.RFC3339), job.ID, imgFileName))
					if !fileExists(imgsImgFileOutput) {

						// Fetch the URL.
						log.Infof("BEGIN Fetching image: %s", imgUrl)
						req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgUrl, nil)
						if err != nil {
							return fmt.Errorf("failed to create HTTP image request for %s : %w", imgUrl, err)
						}
						resp, err := http.DefaultClient.Do(req)

						if err == nil {
							imgContents, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								return fmt.Errorf("failed to read image HTTP body url=%s %w", imgUrl, err)
							}
							imgFileName := path.Base(imgUrl)

							// Save just image
							createdAt := job.EnqueueTime.Local()
							imgsImgFileOutput := filepath.Join(imgsDir, fmt.Sprintf("%s_%s_%s", createdAt.Format(time.RFC3339), job.ID, imgFileName))
							err = os.WriteFile(imgsImgFileOutput, imgContents, 0644)
							if err != nil {
								return fmt.Errorf("failed to write image file =%s %w", imgsImgFileOutput, err)
							}
							// Save image to job directory also
							err = os.WriteFile(filepath.Join(jobDir, imgFileName), imgContents, 0644)
							if err != nil {
								return fmt.Errorf("failed to write image file JSON jobID=%s img=%s %w", job.ID, imgUrl, err)
							}

							resp.Body.Close()
							log.Infof("DONE Fetching image: %s", imgUrl)
						} else {
							log.Errorf("ERROR Fetching image: %s %v", imgUrl, err)
							return err
						}
					} else {
						log.Infof("SKIPPING - File already exists for url: %s : %s", imgUrl, imgsImgFileOutput)
					}

				}
				return nil
			})

		}

		err = g.Wait()
		if err != nil {
			return fmt.Errorf("error fetching images from recent jobs %w", err)
		}
		return nil
		//return render(cmd.OutOrStdout(), format, r)
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
