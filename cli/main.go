// To compile this cli, you must copy all code from insapp-go/src.
// It replace the main.go given in ../src
// This is done automatically if you're using the Dockerfile

package main

import (
	"fmt"
	"log"
	"errors"
	"os"

	"github.com/urfave/cli"
)

func main(){
	app := cli.NewApp();
	app.Name = "Insapp-api-cli"
	app.Usage = "A useful cli to manage insapp api"
	app.Version = "0.0.1"
		app.Authors = []cli.Author{
		cli.Author{
			Name:  "Pitou Games",
			Email: "pitou.games@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 Insapp"

	app.UseShortOptionHandling = true
	app.Commands = []cli.Command{

		cli.Command{
			Name: "association",
			Category: "setup",
			Usage: "Manage associations",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "Create a master association",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name",
							Usage: "Name of the association",
						},
						cli.StringFlag{
							Name: "email",
							Usage: "Email to contact the association",
						},
					},
					Action: func(c *cli.Context) error {
						if c.String("name") != "" && c.String("email") != "" {
							err := AddAssociationCLI(c.String("name"), c.String("email"))
							return err
						} else {
							return errors.New("You must provide a --name and an --email. See association create --help")
						}
					},
				},
			},
		},

		cli.Command{
			Name: "cdn-clean",
			Category: "management",
			Usage: "Clean cdn by removing all images not referenced in database",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "archive, a",
					Usage: "Move files in an archive sub-folder",
				},
				cli.BoolFlag{
					Name: "delete, d",
					Usage: "Delete files forever",
				},
				cli.BoolFlag{
					Name: "list, l",
					Usage: "List all files that will be affected",
				},
			},
			Action: func(c *cli.Context) error {
				var usedImages = GetUsedImages()
				cdnImages, _ := GetImagesNames()
				var toDelete []string
				for _, cdnImage := range cdnImages {
				    delete := true
				    for _, usedImage := range usedImages {
                        if usedImage == cdnImage {
				    		delete = false
							break
                        }
                    }
                    if delete {
                    	toDelete = append(toDelete, cdnImage)
                    }
				}

				if c.Bool("list") {
					for _, imageName := range toDelete {
						fmt.Println(imageName)
					}
				}

				fmt.Println(len(usedImages), " images found in database")
				fmt.Println(len(cdnImages), " images found in cdn")
				fmt.Println(len(toDelete), " images will be affected")

				if c.Bool("archive") {
					fmt.Println("Archiving files...")
					for _, imageName := range toDelete {
						err := ArchiveImage(imageName);
						if err != nil {
							return err
						}
					}
					fmt.Println("Done!")
					return nil
				} else if c.Bool("delete") {
					fmt.Println("Deleting files...")
					for _, imageName := range toDelete {
						err := DeleteImage(imageName);
						if err != nil {
							return err
						}
					}
					fmt.Println("Done!")
					return nil
				} else {
					fmt.Println("To list files affected, use -l")
					fmt.Println("If you're sure to delete these files, use -d but we suggest to archive them with -a")
					return nil
				}
			},
		},

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// brand new created master association
func AddAssociationCLI(name string, email string) error {
	var association Association
	association = GetAssociationFromEmail(email)
	if association.Email != "" {
		return errors.New("This email is already used by " + association.Name)
	}

	fmt.Println("Creating Association:", name, email)
	association.Name = name
	association.Email = email
	res := AddAssociation(association)
	password := GeneratePassword()
	fmt.Println("Association created:", res)

	var user AssociationUser
	user.Association = res.ID
	user.Username = res.Email
	user.Master = true
	user.Password = GetMD5Hash(password)
	AddAssociationUser(user)
	err := SendAssociationEmailSubscription(user.Username, password)
	if err != nil {
		return err
	}
	fmt.Println("An email has been sent to the given address, containing your credential")
	return nil
}

func GetUsedImages() []string {
	var result []string;
	var assos = GetAllAssociation()
    for _, ass := range assos {
    	if ass.Profile != "" {
    		result = append(result, ass.Profile)
    	}
		if ass.Cover != "" {
    		result = append(result, ass.Cover)
		}
    }

	var events = GetEvents()
    for _, event := range events {
		if event.Image!= "" {
			result = append(result, event.Image)
		}
    }

	var posts = GetPosts()
    for _, post := range posts {
		if post.Image!= "" {
			result = append(result, post.Image)
		}
    }

	return result;
}
