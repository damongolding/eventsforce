package serve

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/damongolding/eventsforce/cmd/build"
	"github.com/damongolding/eventsforce/internal/configuration"
	"github.com/damongolding/eventsforce/internal/utils"

	"github.com/jaschaephraim/lrserver"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	config  configuration.Config
	devPort string
)

func init() {
	config = *configuration.NewConfig()
}

// rootCmd represents the base command when called without any subcommands
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start dev server",
	Long:  "Start dev server",

	Run: func(cmd *cobra.Command, args []string) {

		devPort = cmd.Flag("port").Value.String()

		start := time.Now()

		defer func(start time.Time) {
			if err := recover(); err != nil {
				b := utils.Red("Failed in")
				t := utils.BoldRedUnderline(fmt.Sprintf("%.2f", time.Since(start).Seconds()))
				s := utils.Red("seconds")
				l := utils.BoldYellow("😔")

				fmt.Println("\n", err)

				fmt.Println("\n", b, t, s, l)
				os.Exit(1)
			}
		}(start)

		if err := Serve(); err != nil {
			panic(err)
		}
	}}

func watcher() {

	// Create and start LiveReload server
	lr := lrserver.New("ef", lrserver.DefaultPort)
	lr.SetStatusLog(nil)
	lr.SetErrorLog(nil)
	lr.SetStatusLog(nil)

	fmt.Println(utils.SectionMessage("Watching", config.SrcDir, "for changes", "👀"))

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		// ticker for debouncing
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {
					var wg sync.WaitGroup
					wg.Add(1)

					if err := build.Build(false, &wg); err != nil {
						panic(err)
					}
					wg.Wait()

					fmt.Println(utils.SectionMessage("Rebuilt"))
					lr.Reload(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}

			// Wait for 1 sec to pass (if it hasnt already)
			select {
			case <-ticker.C:
				continue
			}
		}
	}()

	// Watch all subdirs
	filepath.Walk(config.SrcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Add a path.to watch
			err = watcher.Add(path)
			if err != nil {
				return err
			}
		}

		return nil

	})

	go lr.ListenAndServe()

	select {}

}

func Serve() error {

	var wg sync.WaitGroup
	wg.Add(1)

	if err := build.Build(false, &wg); err != nil {
		return err
	}

	wg.Wait()

	if err := config.PrintConfig(); err != nil {
		return err
	}

	serverSectionTitle, err := utils.OutputSectionStyling("S", "E", "R", "V", "E")
	if err != nil {
		return err
	}

	fmt.Println(serverSectionTitle)

	go watcher()

	url := fmt.Sprintf("http://localhost:%s", devPort)

	http.Handle("GET /", http.FileServer(http.Dir(config.BuildDir)))
	http.Handle("GET /_assets/", http.FileServer(http.Dir(config.SrcDir)))

	fmt.Println(utils.SectionMessage("Serving on", url))

	err = utils.Openbrowser(url)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(fmt.Sprintf(":%s", devPort), nil)
	if err != nil {
		return err
	}

	return nil
}
