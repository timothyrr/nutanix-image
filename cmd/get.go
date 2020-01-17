package cmd

import (
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"***REMOVED***/nutanix-image/get"
)

// nutanixImageGetCmd represents the get command
var nutanixImageGetCmd = &cobra.Command{
	Use:   "get <image-name>",
	Short: "Download a nutanix image",
	Long: `This command will download a given image by name
using the filepath you specify. If no filepath is specified
it will use a default path of your current working directory.

Example:

  $ nutanix-image get rhel8 --output_dir=/var/images

`,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		requireSliceFlag("endpoint")
		requireStringFlag("username")
		requireStringFlag("password")

        endpoints := viper.GetStringSlice("endpoint")
        for _, endpoint := range endpoints {
        	log.Infof(endpoint)
            checks := []string{"https://", "http://", ":9440"}
            for i := 0; i < len(checks); i++ {
                if strings.Contains(endpoint, checks[i]) {
                	log.Infof(checks[i])
		        	log.Fatal("Invalid nutanix endpoint provided. No need to specify \"https://\" or \"http://\" and/or the port number (9440). Only need to specify the endpoint (e.g. nutanix.example.com)")
		        }
		    }
		    nc := get.NewNutanixClient(
		    	endpoint,
		    	viper.GetString("username"),
		    	viper.GetString("password"),
		    	viper.GetBool("insecure"),
		    	viper.GetBool("debug"),
		    )


		    output_dir := viper.GetString("output_dir")

		    for _, name := range args {
		    	log.Infof("Locating Nutanix image (%s)...", name)
		    	image, err := nc.GetImage(name)
		    	if err != nil {
		    		log.Fatal(err)
		    	}
		        switch num := len(image.Entities); {
		        case num == 0:
		        	log.Infof("No Nutanix images found for: %s", name)
		        case num > 1:
		        	log.Fatalf("Multiple Nutanix images found for %s. Aborting...", name)
		        case num == 1:
		        	i := image.Entities[0]
		        	uuid := i.Metadata.UUID
		        	out_file := fmt.Sprintf("%s/%s.img", output_dir, name)
		        	log.Infof("Found Nutanix image with UUID: %s", uuid)
		        	log.Infof("Downloading Nutanix image file to %s...", out_file)
		        	err = nc.DownloadImage(out_file, uuid, name)
		        	if err != nil {
		        		log.Fatal(err)
		        	}
		        	log.Infof("Nutanix image downloaded")
		        }
		    }
		}
	},
}

func init() {
    path, err := os.Getwd()
    if err != nil {
        log.Println(err)
    }

	rootCmd.AddCommand(nutanixImageGetCmd)
    nutanixImageGetCmd.PersistentFlags().String("output_dir", path, "the directory to save the image to")
	viper.BindPFlag("output_dir", nutanixImageGetCmd.PersistentFlags().Lookup("output_dir"))
}
