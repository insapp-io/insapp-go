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

	fmt.Fprintln(w, "<h2>C'est quoi ?</h2>")

	fmt.Fprintln(w, `<p>
Insapp est une application promouvant l'associatif au sein de l'INSA de Rennes.
L'application a été initialement imaginée et développée par Antoine Crochet et Florent Thomas-Morel en 2016, en collaboration avec l'AEIR (notamment avec Théau Jubin et Antoine Tulasne) et le CRI de l'INSA de Rennes.
Thomas Bouvier, Louis-Sinan Cappoen, et Guillaume Courtet ont ensuite rejoint l'équipe de développement.
		</p>`)

	fmt.Fprintln(w, "<h2>Données</h2>")

	fmt.Fprintln(w, `<p>
Les données transmises à l'application sont matériellement stockées sur les serveurs du CRI.
Vos identifiants INSA, utilisés lors de l'inscription ne sont pas stockés pour Insapp pour des raisons de sécurité et de vie privée.
À aucun moment nous ne stockons, diffusons ou avons accès à vos identifiants INSA.
Vous n'êtes en aucun cas obligé de renseigner votre identité dans la page profil de l'application, qui n'est là qu'à titre facultatif.
		</p>`)

	fmt.Fprintln(w, "<h2>Associations</h2>")

	fmt.Fprintln(w, `<p>
Les associations ont accès au nombre de "like", le nombre de participants pour leurs évènements, ainsi qu'au contenu des commentaires postés sur leurs posts.
Ces commentaires restent anonymes et les associations se réservent le droit de supprimer des commentaires.
Les notifications peuvent être activées ou désactivées dans les réglages du téléphone.
		</p>`)

	fmt.Fprintln(w, "<h2>Testeurs</h2>")

	fmt.Fprintln(w, `<p>
Insapp a été testé pendant 3 semaines par des beta-testeurs volontaires issus de différents départements de l'INSA. Nous tenons à remercier pour leurs retours :
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
Tanguy Le Quéré</br>
Laurent Quénach de Quivillic</br>
Florian Arnoud</br>
Anthony Sharpe</br>
Sebastien Turpin</br>
Luc Geffrault</br>
Agathe Duboue</br>
Alexis Brard</br>
Timothé Frignac</br>
</br>
		</p>`)

	fmt.Fprintln(w, "<h2>Sources</h2>")

	fmt.Fprintln(w, `<p>
Le code source d'Insapp est ouvert et libre de droits  :
</br>
</br>
https://github.com/tomatrocho/insapp-web (forké de https://github.com/fthomasmorel/insapp-web)</br>
https://github.com/tomatrocho/insapp-go (forké de https://github.com/fthomasmorel/insapp-go)</br>
https://github.com/fthomasmorel/insapp-iOS</br>
https://github.com/tomatrocho/insapp-android</br>
</br>
		</p>`)

	fmt.Fprintln(w, "</body>")
}

func Legal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<body style='font-family: \"Arial\", Arial, sans-serif; text-align: justify;'>")
	fmt.Fprintln(w, "<h1>Insapp</h1>")

	fmt.Fprintln(w, `<p>
Les données transmises à l’application sont matériellement stockées sur les serveurs loués par l'Amicale des élèves. À aucun moment Insapp ne stocke, diffuse ou a accès aux identifiants INSA des utilisateurs, pour des raisons de sécurité et de confidentialité.
		</p>`)

	fmt.Fprintln(w, `<p>
L'Utilisateur s'assure de garder son mot de passe secret. Toute divulgation du mot de passe, quelle que soit sa forme, est interdite. L'Utilisateur assume les risques liés à l'utilisation de son identifiant et mot de passe. Insapp décline toute responsabilité. Une garantie optimale de la sécurité et de la confidentialité des données transmises n'est pas assurée par Insapp. Toutefois, Insapp s'engage à mettre en œuvre tous les moyens nécessaires afin de garantir au mieux l’intégrité des données.
		</p>`)

	fmt.Fprintln(w, `<p>
Les sources des informations diffusées sur Insapp sont réputées fiables. Toutefois, Insapp se réserve la faculté d'une non-garantie de la fiabilité des sources. Les informations données sur Insapp le sont à titre purement informatif. Ainsi, l'Utilisateur assume seul l'entière responsabilité de l'utilisation des informations et contenus.
		</p>`)

	fmt.Fprintln(w, `<p>
L’application permet aux membres de publier des commentaires. Dans ses publications, l'Utilisateur s’engage à respecter la charte informatique de l’INSA de Rennes et les règles de droit en vigueur. Insapp exerce une modération a posteriori sur les publications et se réserve le droit de refuser leur mise en ligne, sans avoir à s’en justifier auprès du membre. L’application permet également à l’utilisateur, s’il le souhaite, de compléter son profil.
		</p>`)

	fmt.Fprintln(w, `<p>
Insapp requiert l’accès à la caméra lors du scan du code barre de la carte amicaliste. Cette donnée est stockée localement, et n’est en aucun cas transmise à un tiers.
		</p>`)

	fmt.Fprintln(w, `<p>
Insapp décline toutes responsabilités en cas de refus ou de dysfonctionnement technique lors de la présentation du code barre de la carte amicaliste.
		</p>`)

	fmt.Fprintln(w, "</body>")
}
