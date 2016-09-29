package main

import (
	"fmt"
	"net/http"
)

// Index is just a test actually
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Insapp REST API - v.0.1")
}



func Credit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<body style='font-family: \"Arial\", Arial, sans-serif; text-align: justify;'>")
	fmt.Fprintln(w, "<h1>Insapp</h1>")

	fmt.Fprintln(w, "<h2>Quoi ?</h2>")

	fmt.Fprintln(w, `<p>
Insapp est une application qui a pour but de promouvoir l'associatif au sein de l'INSA Rennes.
L'application a été initialement imaginée et développée par Antoine Crochet et Florent Thomas-Morel en 2016.
Nous avons travaillé de pair avec l'AEIR et notamment Théau Jubin et Antoine Tulasne ainsi qu'avec le CRI de l'INSA Rennes.
		</p>`)

	fmt.Fprintln(w, "<h2>Données</h2>")

	fmt.Fprintln(w, `<p>
Les données transmises à l'application sont matériellement stockées sur les serveurs du CRI.
Vos identifiants INSA, utilisés lors de l'inscription ne sont pas stockés pour Insapp pour des raisons de sécurité et de vie privée.
À aucun moment nous ne stockons, diffusons ou avons accès à vos identifiants INSA.
Vous n'êtes en aucun cas obligé de renseigner votre identité dans la page profil de l'application qui n'est là qu'à titre facultatif.
		</p>`)

	fmt.Fprintln(w, "<h2>Associations</h2>")

	fmt.Fprintln(w, `<p>
Les associations ont accès au nombre de "like", le nombre de participants pour leurs évènements, ainsi qu'au contenu des commentaires postés sur leurs posts.
Ces commentaires restent anonymes et les associations se réservent le droit de supprimer des commentaires.
Les notifications peuvent être activées ou désactivées dans les réglages du téléphone.
		</p>`)

	fmt.Fprintln(w, "</body>")
}
