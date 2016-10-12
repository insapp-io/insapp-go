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


	fmt.Fprintln(w, "<h2>Testeurs</h2>")

	fmt.Fprintln(w, `<p>
Insapp a été testé pendant 3 semaines par des beta-testeurs volontaires issuent de differents départements de l'INSA. Nous tenons à remercier :
</br>
</br>
Nans Préjean</br>
Hugo David</br>
Valentin Marc</br>
Alex Gravier</br>
Mathieu Cassard</br>
Célestin Bodet</br>
Etienne Rebout</br>
Rémy Garcia</br>
Jean-Baptiste Nou</br>
Thomas Bouvier</br>
Tanguy Le Quéré</br>
Laurent Quénach de Quivillic</br>
Florian Arnoud</br>
Anthony Sharpe</br>
Sebastien Turpin</br>
Luc Geffrault</br>
Agathe Duboue</br>
</br>
		</p>`)

	fmt.Fprintln(w, "</body>")
}




func Legal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<body style='font-family: \"Arial\", Arial, sans-serif; text-align: justify;'>")
	fmt.Fprintln(w, "<h1>Insapp</h1>")

	fmt.Fprintln(w, `<p>
Les sources des informations diffusées sur INSAPP sont réputées fiables. Toutefois, INSAPP se réserve la faculté d'une non-garantie de la fiabilité des sources. Les informations données sur le site le sont à titre purement informatif. Ainsi, l'Utilisateur assume seul l'entière responsabilité de l'utilisation des informations et contenus du présent site.
L'Utilisateur s'assure de garder son mot de passe secret. Toute divulgation du mot de passe, quelle que soit sa forme, est interdite.
		</p>`)

	fmt.Fprintln(w, `<p>
L'Utilisateur assume les risques liés à l'utilisation de son identifiant et mot de passe. Le site décline toute responsabilité.
Une garantie optimale de la sécurité et de la confidentialité des données transmises n'est pas assurée par INSAPP. Toutefois, le site s'engage à mettre en œuvre tous les moyens nécessaires afin de garantir au mieux la sécurité et la confidentialité des données.
		</p>`)

	fmt.Fprintln(w, `<p>
Le site permet aux membres de publier des commentaires.
Dans ses publications, l'Utilisateur s’engage à respecter la charte informatique de l’INSA Rennes et les règles de droit en vigueur.
Le site exerce une modération a posteriori sur les publications et se réserve le droit de refuser leur mise en ligne, sans avoir à s’en justifier auprès du membre.
		</p>`)

	fmt.Fprintln(w, "</body>")
}
