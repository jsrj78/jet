(ns app.views
  (:require [reagent.core :as reagent]
            [re-frame.core :as rf]
            [clojure.string :as str]
            [goog.events :as events]
            [cljs.pprint :refer [pprint]]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref rf/subscribe))
(def >evt rf/dispatch)

(defn client-xy [evt]
  [(.-clientX evt) (.-clientY evt)])

(defn drag-move-fn [id *state]
  (fn [evt]
    (let [[ox oy]       (:pos @*state)
          [cx cy :as c] (client-xy evt)]
      (swap! *state assoc :pos c)
      (>evt [:move-gadget id (- cx ox) (- cy oy)])
      (.preventDefault evt)))) ; avoid text selection, at least in Chrome

(defn drag-end-fn [move-fn *state]
  (fn [evt]
    (events/unlisten js/window "mousemove" move-fn)
    (events/unlisten js/window "mouseup" (:end @*state))))

(defn drag-start [id x y evt]
  (let [*state  (atom {:pos (client-xy evt)})
        move-fn (drag-move-fn id *state)
        done-fn (drag-end-fn move-fn *state)]
    (swap! *state assoc :end done-fn)
    (events/listen js/window "mousemove" move-fn)
    (events/listen js/window "mouseup" done-fn)
    (>evt [:select-gadget id])
    (.stopPropagation evt)))

(defn get-dom-width [elt]
  ; getBBox will be found via the "svg.js" externs file available here:
  ;  https://github.com/google/closure-compiler/blob/master/contrib/externs/
  (-> elt reagent/dom-node .getBBox .-width int))

(defn adjust-label-width [id x y label]
  ; see https://stackoverflow.com/questions/27602592/reagent-component-did-mount
  ; inlined, this adjusts the enclosing rect's width to the text after render
  [^{:component-did-mount #(>evt [:set-label-width id (get-dom-width %)])}
    #(do [:text.obj {:x x :y y} label])])

(defn obj-as-svg [id x y w h label]
  [:g
    [:rect.obj {:x x :y y :width w :height h}]
    (adjust-label-width id (+ x 5) (+ y 15) label)])

(defn msg-as-svg [id x y w h label]
  [:g
    [:rect.msg {:x (- x 10) :y y :width (+ w 3) :height h :rx (/ h 2)
                :on-mouse-down #(>evt [:to-engine id])}]
    (adjust-label-width id (- x 3) (+ y 15) label)])

(defn bang-as-svg [id x y]
  [:circle.bang {:cx (+ x 2.5) :cy (+ y 10) :r 10
                 :on-mouse-down #(>evt [:to-engine id])}])

(defn toggle-as-svg [id x y]
  [:rect.toggle {:x (- x 7.5) :y y :width 20 :height 20}])

(defn gadget-as-svg [id obj]
  (let [[gbox ic oc] (<sub [:gadget-coords id])
        [x y w h]    gbox]
    ^{:key gbox}
     [:g.draggable {:on-mouse-down #(drag-start id x y %)}
      (map (fn [[cx cy :as xy]]
             ^{:key xy}
              [:circle {:cx cx :cy cy :r 3}]) (concat ic oc))
      (if (= (nth obj 2) :obj)
        (obj-as-svg id x y w h (<sub [:gadget-name id]))
        (case (nth obj 3)
          :bang   (bang-as-svg id x y)
          :toggle (toggle-as-svg id x y)
          :msg    (msg-as-svg id x y w h (<sub [:gadget-name id]))))]))

(defn wire-path [[x1 y1] [x2 y2]] ;; either straight line or cubic bezier
; (str/join " " ["M" x1 y1 "L" x2 y2])
  (str/join " " ["M" x1 y1 "C" x1 (+ y1 25) x2 (- y2 25) x2 y2]))

(defn wire-as-svg [[s-id s-outlet d-id d-inlet :as wire]]
  (let [[_ _ s-outs] (<sub [:gadget-coords s-id])
        [_ d-ins _]  (<sub [:gadget-coords d-id])]
    ^{:key wire}
     [:path.wire {:d (wire-path (nth s-outs s-outlet)
                                (nth d-ins d-inlet))}]))

(defn design-as-svg []
  (let [objs  (<sub [:gadgets])
        wires (<sub [:wires])]
    [:svg {:width "100%" :height 400
           :on-mouse-down #(>evt [:select-gadget nil])}
      ; can't leave reactive refs in a lazy sequences
      (doall (map-indexed gadget-as-svg objs))
      (doall (map wire-as-svg wires))]))

(defn main-menu []
  [:div#menu.custom-wrapper.pure-g
    [:div.pure-u-1.pure-u-md-1-3
      [:div.pure-menu
        [:a.pure-menu-heading.custom-brand {:href "#"} "Brand"]
        [:a#toggle.custom-toggle {:href "#"} [:s.bar] [:s.bar]]]]
    [:div.pure-u-1.pure-u-md-1-3
      [:div.pure-menu.pure-menu-horizontal.custom-can-transform
        [:ul.pure-menu-list
          [:li.pure-menu-item [:a.pure-menu-link {:href "#"} "Home"]]
          [:li.pure-menu-item [:a.pure-menu-link {:href "#"} "About"]]
          [:li.pure-menu-item [:a.pure-menu-link {:href "#"} "Blah"]]]]]
    [:div.pure-u-1.pure-u-md-1-3
      [:div.pure-menu.pure-menu-horizontal.custom-menu-3.custom-can-transform
        [:ul.pure-menu-list
          [:li.pure-menu-item [:a.pure-menu-link {:href "#"} "Foo"]]
          [:li.pure-menu-item [:a.pure-menu-link {:href "#"} "Bar"]]]]]])

(defn design-inspector []
  [:p (count (<sub [:gadgets])) " gadgets, "
      (count (<sub [:wires])) " wires"])

(defn gadget-inspector [obj]
  [:div
    [:pre (str obj)]])

(defn main-content []
  [:div
    [:div.pure-g.pure-u-3-5
     [:div#content
      [design-as-svg]
      ;[:pre [:small (pr-str @re-frame.db/app-db)]]
      [:pre [:small (with-out-str (pprint @re-frame.db/app-db))]]]]
    [:div.pure-g.pure-u-2-5
     [:div#sidebar
      (if-let [obj (<sub [:curr-gadget])]
        [gadget-inspector obj]
        [design-inspector])]]])

(defn app-page []
  [:div
    #_[main-menu]
    [main-content]])
