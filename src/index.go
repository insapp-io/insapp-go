package main

import (
	"fmt"
	"net/http"
)

// Index is just a test actually
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Insapp REST API - v1.0")
}

func HowToPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<body style='font-family: \"Arial\", Arial, sans-serif; text-align: justify;'>")
	fmt.Fprintln(w, "<h1>Insapp</h1>")

	fmt.Fprintln(w, "<h2>Comment publier ?</h2>")

	fmt.Fprintln(w, `<p>
L'adresse de l'interface d'administration est la suivante : insapp.fr/admin.
		</p>`)

	fmt.Fprintln(w, `<p>
Pour pouvoir publier depuis l'interface d'administration, un compte d'association est nécessaire.
Si votre association n'en dispose pas encore, ou si vous avez oublié le mot de passe, n'hésitez pas à envoyer un mail à aeir-insapp@insa-rennes.fr.
Nous répondons sous peu !
		</p>`)

	fmt.Fprintln(w, "</body>")
}

func Credit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<body style='font-family: \"Arial\", Arial, sans-serif; text-align: justify;'>")
	fmt.Fprintln(w, "<h1>Insapp</h1>")

	fmt.Fprintln(w, "<h2>C'est quoi ?</h2>")

	fmt.Fprintln(w, `<p>
Insapp est une application promouvant l'associatif au sein de l'INSA Rennes.
L'application a été initialement imaginée et développée par Antoine Crochet et Florent Thomas-Morel début 2016. Dès octobre 2016, Thomas Bouvier et Guillaume Courtet ont rejoint l'équipe de développement.
Thomas Bouvier, Pierre Duc-Martin ainsi qu'Antoine Pégné maintiennent actuellement le projet.
Insapp est un projet lancé en collaboration avec l'AEIR et le CRI de l'INSA Rennes. Merci à Théau Jubin et Antoine Tulasne (2016), Laura Frouin (2017) et Titouan Le Hir (2018) pour leurs contributions !
		</p>`)

	fmt.Fprintln(w, "<h2>Données</h2>")

	fmt.Fprintln(w, `<p>
Vos identifiants INSA, utilisés lors de l'inscription ne sont pas stockés pour Insapp pour des raisons de sécurité et de vie privée.
À aucun moment nous ne stockons, diffusons ou avons accès à vos identifiants INSA.
Vous n'êtes en aucun cas obligé de renseigner votre identité dans la page profil de l'application, qui n'est là qu'à titre facultatif.
		</p>`)

	fmt.Fprintln(w, "<h2>Associations</h2>")

	fmt.Fprintln(w, `<p>
Les associations ont accès au nombre de "like", le nombre de participants pour leurs évènements, ainsi qu'au contenu des commentaires postés sur leurs posts.
Ces commentaires restent anonymes et les associations se réservent le droit d'en supprimer.
Les notifications peuvent être activées ou désactivées dans les paramètres de l'application.
		</p>`)

	fmt.Fprintln(w, "<h2>Sources</h2>")

	fmt.Fprintln(w, `<p>
Le code source d'Insapp est ouvert et libre de droits. Le code des applications Android (Java & Kotlin) et iOS (Swift), de l'interface web (AngularJS) et de l'API (Golang) est accessible sur Github, aux adresses suivantes :
</br>
</br>
https://github.com/thomas-bouvier/insapp-server.git</br>
https://github.com/thomas-bouvier/insapp-android.git</br>
https://github.com/RobAddict/insapp-iOS.git</br>
</br>
		</p>`)

	fmt.Fprintln(w, "</body>")
}

func Legal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<body style='font-family: \"Arial\", Arial, sans-serif; text-align: justify;'>")
	fmt.Fprintln(w, "<h1>Insapp</h1>")

	fmt.Fprintln(w, `<p>
Les données transmises à l’application sont matériellement stockées sur les serveurs loués par l'AEIR. À aucun moment Insapp ne stocke, diffuse ou a accès aux identifiants INSA des utilisateurs, pour des raisons de sécurité et de confidentialité.
		</p>`)

	fmt.Fprintln(w, `<p>
L'Utilisateur s'assure de garder son mot de passe secret. Toute divulgation du mot de passe, quelle que soit sa forme, est interdite. L'Utilisateur assume les risques liés à l'utilisation de son identifiant et mot de passe. Insapp décline toute responsabilité. Une garantie optimale de la sécurité et de la confidentialité des données transmises n'est pas assurée par Insapp. Toutefois, Insapp s'engage à mettre en œuvre tous les moyens nécessaires afin de garantir au mieux l’intégrité des données.
		</p>`)

	fmt.Fprintln(w, `<p>
Les sources des informations diffusées sur Insapp sont réputées fiables. Toutefois, Insapp se réserve la faculté d'une non-garantie de la fiabilité des sources. Les informations données sur Insapp le sont à titre purement informatif. Ainsi, l'Utilisateur assume seul l'entière responsabilité de l'utilisation des informations et contenus.
		</p>`)

	fmt.Fprintln(w, `<p>
L’application permet aux membres de publier des commentaires. Dans ses publications, l'Utilisateur s’engage à respecter la charte informatique de l’INSA Rennes et les règles de droit en vigueur. Insapp exerce une modération a posteriori sur les publications et se réserve le droit de refuser leur mise en ligne, sans avoir à s’en justifier auprès du membre. L’application permet également à l’utilisateur, s’il le souhaite, de compléter son profil.
		</p>`)

	fmt.Fprintln(w, `<p>
Insapp requiert l’accès à la caméra lors du scan du code barre de la carte amicaliste. Cette donnée est stockée localement, et n’est en aucun cas transmise à un tiers.
		</p>`)

	fmt.Fprintln(w, `<p>
Insapp décline toutes responsabilités en cas de refus ou de dysfonctionnement technique lors de la présentation du code barre de la carte amicaliste.
		</p>`)

	fmt.Fprintln(w, "</body>")
}
