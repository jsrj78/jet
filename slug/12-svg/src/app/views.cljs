(ns app.views
  (:require [re-frame.core]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn as-svg [[_ x y] obj]
  ^{:key (+ (* x 10000) y)}
  [:rect.obj {:x x :y y :width 30 :height 20}])

(defn svg-component []
  [:svg {:id "canvas"
         :width "400" :height "200"
         :style {:outline "1px solid black"}}
        [:circle.test {:cx 180 :cy 80 :r 30}]
        [:rect.test {:x 250 :y 50 :width 50 :height 30}]
        (map as-svg (filter #(= (first %) :obj) (<sub [:design])))])

(defn main-panel []
  (let [name (<sub [:name])] 
    ;; this displays on the JavaScript console
    ;(.log js/console "db:" (pr-str @re-frame.db/app-db))

   [:div#main

    [:div.pure-g
     [:div.pure-u-1
      [:h3 [:b "Hello " name]]
      [:p
       [:input {:type "text"
                :value name
                :on-change #(>evt [:change-name (.. % -target -value)])}]]]]

    [:div.pure-g
     [:div.pure-u-1
      [svg-component]]]

    [:div.pure-g
     [:div.pure-u-1
      ;; this displays as <pre> text in a verbose format
      ;[:pre (with-out-str (cljs.pprint/pprint @re-frame.db/app-db))]
      ;; this displays as text on a single line
      [:p (pr-str @re-frame.db/app-db)]]]]))
