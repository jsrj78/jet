(ns app.views
  (:require [re-frame.core]
            [clojure.string :as s]
            [goog.events :as ev]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-name [[_ _ _ & cmd]]
  (subs (s/join " " cmd) 1))

(defn obj-id [[_ x y]]
  (str x "," y))

(defn client-xy [evt]
  [(.-clientX evt) (.-clientY evt)])

(defn drag-move-fn [id state]
  (fn [evt]
    (let [[ox oy]       (:pos @state)
          [cx cy :as c] (client-xy evt)]
      (swap! state assoc :pos c)
      (>evt [:move-gadget id (- cx ox) (- cy oy)]))))

(defn drag-end-fn [id move-fn state]
  (fn [evt]
    (>evt [:select-gadget id])
    (ev/unlisten js/window "mousemove" move-fn)
    (ev/unlisten js/window "mouseup" (:end @state))))

(defn drag-start [x y evt]
  (let [state   (atom {:pos (client-xy evt)})
        id      (js/parseInt (.-id (.-target evt)))
        move-fn (drag-move-fn id state)
        done-fn (drag-end-fn id move-fn state)]
    (swap! state assoc :end done-fn)
    (ev/listen js/window "mousemove" move-fn)
    (ev/listen js/window "mouseup" done-fn)
    #_(.stopPropagation evt)))

(defn get-dom-width [elt]
  (.. (reagent.core/dom-node elt) getBBox -width))

(defn obj-as-svg [id [_ x y & cmd :as obj]]
  ; see https://stackoverflow.com/questions/27602592/reagent-component-did-mount
  (let [new-width #(>evt [:set-label-width id (+ (get-dom-width %) 10)])
        did-mount (with-meta identity {:component-did-mount new-width})]
    ^{:key id}
    [:g.draggable {:on-mouse-down #(drag-start x y %)}
      [:rect.obj {:id id :x x :y y :width (<sub [:label-width id]) :height 20}]
      [did-mount
        [:text.obj {:x (+ x 5) :y (+ y 15)} (obj-name obj)]]]))

(defn wire-id [wire]
  (s/join ":" wire))

(defn wire-path [x1 y1 x2 y2] ;; either straight line or cubic bezier
; (s/join " " ["M" x1 y1 "L" x2 y2])
  (s/join " " ["M" x1 y1 "C" x1 (+ y1 25) x2 (- y2 25) x2 y2]))

(defn wire-as-svg [[src-pos src-out dst-pos dst-in :as wire]]
  (let [[_ sx sy] (<sub [:gadget-num src-pos])
        [_ dx dy] (<sub [:gadget-num dst-pos])
        src-width (<sub [:label-width src-pos])
        dst-width (<sub [:label-width dst-pos])]
    ^{:key (wire-id wire)}
    [:path.wire {:d (wire-path (+ sx (* src-width src-out)) (+ sy 20)
                               (+ dx (* dst-width dst-in))  (+ dy 0))}])) 

(defn design-as-svg []
  (let [objs   (<sub [:gadgets])
        wires  (<sub [:wires])]
    [:svg {:width "100%" :height 500
           :on-mouse-down #(>evt [:select-gadget nil])}
      ; can't leave reactive refs in a lazy sequences
      (doall (map-indexed obj-as-svg objs))
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

(defn gadget-inspector [id]
  [:div
    [:p "gadget: " id]
    [:pre (str (<sub [:gadget-num id]))]])

(defn main-content []
  [:div
    [:div.pure-g.pure-u-3-5
     [:div#content
      [design-as-svg]
      ;[:pre [:small (pr-str @re-frame.db/app-db)]]
      [:pre [:small (with-out-str (cljs.pprint/pprint @re-frame.db/app-db))]]]]
    [:div.pure-g.pure-u-2-5
     [:div#sidebar
      (let [id (<sub [:current-gadget])] 
        (if id
          [gadget-inspector id]
          [design-inspector]))]]])

(defn app-page []
  [:div
    #_[main-menu]
    [main-content]])
