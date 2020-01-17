package cmd

import (
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timothyrr/nutanix-image/get"
	"github.com/timothyrr/nutanix-image/create"
)

// nutanixImageUploadCmd represents the create command
var nutanixImageUploadCmd = &cobra.Command{
	Use:   "create <image-name> --source=<source-file>",
	Short: "Create and upload a nutanix image",
	Long: `This command will create and upload a given image
by name using the filepath you specify with the --source flag.

Example:

  $ nutanix-image create rhel8 --source=/var/images/rhel8.img

`,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		requireSliceFlag("endpoint")
		requireStringFlag("username")
		requireStringFlag("password")
		requireStringFlag("source")

        endpoints := viper.GetStringSlice("endpoint")
        for _, endpoint := range endpoints {
            checks := []string{"https://", "http://", ":9440"}
            for i := 0; i < len(checks); i++ {
                if strings.Contains(endpoint, checks[i]) {
		        	log.Fatal("Invalid nutanix endpoint provided. No need to specify \"https://\" or \"http://\" and/or the port number (9440). Only need to specify the endpoint (e.g. nutanix.example.com)")
		        }
		    }
		    ndc := get.NewNutanixClient(
		    	endpoint,
		    	viper.GetString("username"),
		    	viper.GetString("password"),
		    	viper.GetBool("insecure"),
		    	viper.GetBool("debug"),
		    )

		    nuc := create.NewNutanixClient(
		    	endpoint,
		    	viper.GetString("username"),
		    	viper.GetString("password"),
		    	viper.GetBool("insecure"),
		    	viper.GetBool("debug"),
		    )

		    source := viper.GetString("source")

		    for _, name := range args {
		    	image, err := ndc.GetImage(name)
		    	if err != nil {
		    		log.Fatal(err)
		    	}
		        switch num := len(image.Entities); {
		        case num == 0:
		        	log.Infof("Creating Nutanix image (%s) in %s...", name, endpoint)
		        	uuid, err := nuc.CreateImage(name)
		        	if err != nil {
		        		log.Fatal(err)
		        	}
		        	log.Infof("Uploading Nutanix image file (%s) to %s...", source, endpoint)
		        	err = nuc.UploadImage(uuid.(string), source)
		        	if err != nil {
		        		log.Fatal(err)
		        	}
		        	log.Infof("Nutanix image uploaded (UUID: %s) to %s", uuid, endpoint)
		        case num > 1:
		        	log.Fatalf("Multiple Nutanix images already exist for %s in %s. Aborting...", name, endpoint)
		        case num == 1:
		        	i := image.Entities[0]
		        	uuid := i.Metadata.UUID
		        	log.Fatalf("Existing Nutanix image found for \"%s\" (UUID: %s) in %s. Aborting...", name, uuid, endpoint)
		        }
		    }
		}
	},
}

func init() {
	rootCmd.AddCommand(nutanixImageUploadCmd)
    nutanixImageUploadCmd.PersistentFlags().String("source", "", "the source image to upload")
	viper.BindPFlag("source", nutanixImageUploadCmd.PersistentFlags().Lookup("source"))
}
