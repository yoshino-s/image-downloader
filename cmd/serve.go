package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/yoshino-s/image-downloader/handler/connect"
	"github.com/yoshino-s/image-downloader/handler/http"
	"github.com/yoshino-s/image-downloader/s3"
)

func init() {
	httpApp.Configuration().Register(serveCmd.Flags())
	s3App.Configuration().Register(serveCmd.Flags())

	rootCmd.AddCommand(serveCmd)
}

var (
	s3App = s3.New()

	imageService = connect.NewImageService()

	httpApp = http.New()

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: `Serve runs the HTTP and gRPC server.`,
		Run: func(cmd *cobra.Command, args []string) {

			app.Append(httpApp)
			app.Append(s3App)
			imageService.SetS3(s3App)

			httpApp.SetImageHandler(imageService)

			app.Go(context.Background())
		},
	}
)
