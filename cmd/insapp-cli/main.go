package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	insapp "github.com/thomas-bouvier/insapp-go"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "insapp-cli"
	app.Usage = "A useful CLI to manage Insapp api"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Pitou Games",
			Email: "pitou.games@gmail.com",
		},
	}
	app.Copyright = "(c) 2019 Insapp"
	//app.UseShortOptionHandling = true
	app.Commands = []cli.Command{

		cli.Command{
			Name:     "association",
			Category: "setup",
			Usage:    "Manage associations",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "Create a master association",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Name of the association",
						},
						cli.StringFlag{
							Name:  "email",
							Usage: "Email to contact the association",
						},
					},
					Action: func(c *cli.Context) error {
						if c.String("name") != "" && c.String("email") != "" {
							err := AddAssociationCLI(c.String("name"), c.String("email"))
							return err
						}
						return errors.New("You must provide a --name and an --email. See association create --help")
					},
				},
				{
					Name:  "update",
					Usage: "Update all associations' profile pictures",
					Action: func(c *cli.Context) error {
						return UpdateAssociationsCLI()
					},
				},
			},
		},

		cli.Command{
			Name:     "cdn",
			Category: "management",
			Usage:    "Manage insapp cdn",
			Subcommands: []cli.Command{
				{
					Name:  "clean",
					Usage: "Clean unused images in insapp-cdn folder",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "archive, a",
							Usage: "Move files in an archive sub-folder",
						},
						cli.BoolFlag{
							Name:  "delete, d",
							Usage: "Delete files forever",
						},
						cli.BoolFlag{
							Name:  "list, l",
							Usage: "List all files that will be affected",
						},
					},
					Action: func(c *cli.Context) error {
						var usedImages = GetUsedImages()
						cdnImages, _ := insapp.GetImagesNames()
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
							fmt.Println("Are you sure to archive these files ? (Y/n)")
							if askForConfirmation("Are you sure to archive these files ? (Y/n)") {
								fmt.Println("Archiving files...")
								for _, imageName := range toDelete {
									err := insapp.ArchiveImage(imageName)
									if err != nil {
										return err
									}
								}
								fmt.Println("Done!")
							}
							return nil
						} else if c.Bool("delete") {
							fmt.Println("Are you sure to delete these files ? (Y/n)")
							if askForConfirmation("Are you sure to delete these files ? (Y/n)") {
								fmt.Println("Deleting files...")
								for _, imageName := range toDelete {
									err := insapp.DeleteImage(imageName)
									if err != nil {
										return err
									}
								}
								fmt.Println("Done!")
							}
							return nil
						} else {
							fmt.Println("To list files affected, use -l")
							fmt.Println("If you're sure to delete these files, use -d but we suggest to archive them with -a")
							return nil
						}
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Comes from https://gist.github.com/albrow/5882501
func askForConfirmation(message string) bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if insapp.ContainsString(okayResponses, response) {
		return true
	} else if insapp.ContainsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println(message)
		return askForConfirmation(message)
	}
}

// AddAssociationCLI create a brand new master association
func AddAssociationCLI(name string, email string) error {
	var association insapp.Association
	isValidEmail := insapp.VerifyEmail(email)
	if !isValidEmail {
		return errors.New("This email is already used")
	}

	fmt.Println("Creating Association:", name, email)
	association.Name = name
	association.Email = email
	res := insapp.AddAssociation(association)
	password := insapp.GeneratePassword()
	fmt.Println("Association created:", res)

	var user insapp.AssociationUser
	user.Association = res.ID
	user.Username = res.Email
	user.Master = true
	user.Password = insapp.GetMD5Hash(password)
	insapp.AddAssociationUser(user)
	err := insapp.SendAssociationEmailSubscription(user.Username, password)
	if err != nil {
		return err
	}
	fmt.Println("An email has been sent to the given address, containing your credential")
	return nil
}

// UpdateAssociationsCLI update all associations
func UpdateAssociationsCLI() error {
	var assos = insapp.GetAllAssociations()
	for _, ass := range assos {
		// Migrate profile picture
		if ass.ProfileUploaded == "" && ass.Profile != "" {
			ass.ProfileUploaded = ass.Profile
			ass.Profile = ""
			insapp.UpdateAssociation(ass.ID, ass)
		} else {
			insapp.UpdateAssociation(ass.ID, ass)
		}
	}
	return nil
}

// GetUsedImages return an array of all images file name found in db
func GetUsedImages() []string {
	var result []string
	var assos = insapp.GetAllAssociations()
	for _, ass := range assos {
		if ass.Profile != "" {
			result = append(result, ass.Profile)
		}
		if ass.ProfileUploaded != "" {
			result = append(result, ass.ProfileUploaded)
		}
		if ass.Cover != "" {
			result = append(result, ass.Cover)
		}
	}

	var events = insapp.GetEvents()
	for _, event := range events {
		if event.Image != "" {
			result = append(result, event.Image)
		}
	}

	var posts = insapp.GetPosts()
	for _, post := range posts {
		if post.Image != "" {
			result = append(result, post.Image)
		}
	}

	return result
}
