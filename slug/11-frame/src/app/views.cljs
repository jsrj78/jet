(ns app.views
  (:require [re-frame.core :as rf]))

; https://github.com/Day8/re-frame/blob/master/src/re_frame/subs.cljc#L67-L107
(def <sub 
  (comp deref re-frame.core/subscribe))

(defn main-panel []
  (let [name (<sub [:name])] 
    ;; this displays on the JavaScript console
    ;(.log js/console "db:" (pr-str @re-frame.db/app-db))

    [:div
      [:h1 "Hello " name]
      [:input {:type "text"
               :value name
               :on-change #(rf/dispatch [:change-name (.. % -target -value)])}]

      ;; this displays as <pre> text in a verbose format
      ;[:pre (with-out-str (cljs.pprint/pprint @re-frame.db/app-db))]

      ;; this displays as text on a single line
      [:p (pr-str @re-frame.db/app-db)]]))
