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
			Name: "root-association",
			Category: "setup",
			Usage: "Create a root association",
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
					return errors.New("You must provide a name and an email. See init-association --help")
				}
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
			},
			Action: func(c *cli.Context) error {
				if c.Bool("archive") {
					// TODO : archive
					fmt.Println("Archiving files...")
					fmt.Println("TODO")
					return errors.New("Not implemented")
				} else {
					// TODO : delete
					fmt.Println("Deleting files...")
					fmt.Println("TODO")
					return errors.New("Not implemented")
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
	association = GetAssociationEmail(email)
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
