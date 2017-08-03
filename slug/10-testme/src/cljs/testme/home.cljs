(ns testme.home
  (:require
    [reagent.core :as r]))


(def names [
"Guillermo Wynn  "
"Darlena Paulin  "
"Mara Edelstein  "
"Hillary Minogue  "
"Eleonor Shear  "
"Wynona Harling  "
"Beth Pires  "
"Nida Lapeyrouse  "
"Ricarda Liska  "
"Synthia Dvorak  "
"Royce Murawski  "
"Nerissa Yocom  "
"Tonisha Cardon  "
"Rosalind Wigington  "
"Shaunta Stenson  "
"Anette Emanuel  "
"Oma Pullum  "
"Moses Hornberger  "
"Yvette Fenstermaker  "
"Sung Summa  "
"Katelynn Bradshaw  "
"Robena Manigo  "
"Dong Olivarria  "
"Margarett Oakman  "
"Roman Sabourin  "
"Ayanna Vito  "
"Shellie Cerrato  "
"Reuben Lash  "
"Amee Horak  "
"Jerold Sisemore  "
"Mellissa Lobel  "
"Kristina Gower  "
"Collen Mccreary  "
"Jamie Drinnon  "
"Wanetta Dismukes  "
"Chu Kiel  "
"Lionel Jaeger  "
"Felica Rollo  "
"Pablo Segers  "
"Tawna Inabinet  "
"Jaleesa Portner  "
"Myrna Ridinger  "
"Annette Lawyer  "
"Mike Shue  "
"Min Angus  "
"Blondell Boise  "
"Laine Shepley  "
"Eddie Cadwallader  "
"Parker Band  "
"Jana Bryer  "
])


(defn display-names
  "Displays APIs that match the search string"
  [search-string]
    [:div {:class "container-fluid"}
     [:center
      (for [person-name names]
        ;; tricky regular expression to see if the search string matches the name
        (if (or (re-find (re-pattern (str "(?i)" search-string)) person-name)
                (= "" search-string))
          ^{ :key (.indexOf names person-name)}
            [:div {:class "col-sm-4"}
            [:div {:class "panel panel-default"}
              [:div {:class "panel-heading"} person-name]]]))]])
           
    

(defn show-home
  []
  ;; First, display banner and search text box.
  (let [search-string (r/atom "")]
    (fn []
      [:div {:class "jumbotron " }
        [:div {:class "container-fluid"}
            [:div {:class "row"}]
              [:div {:class "col-sm-2"}]
              [:div {:class "col-sm-8"}
                [:h2 [:center "Fun names" ]]
                [:p {:class "center"} "Filter names" ]
              [:div {:class "row"}]
              [:div {:class "col-sm-2"}]
              [:div {:class "col-sm-8"}
                [:input.form-control {
                  :type "text"
                  :placeholder "Filter names"
                  :value @search-string
                  :on-change #(reset! search-string (-> % .-target .-value))
                }]
           ]
          ]
      ]
    [:div {:class "page-header"}]
    [display-names @search-string]])))
