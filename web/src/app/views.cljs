(ns app.views
  (:require [re-frame.core]
            [clojure.string :as str]
            [goog.events :as ev]))

; see https://lambdaisland.com/blog/11-02-2017-re-frame-form-1-subscriptions
(def <sub (comp deref re-frame.core/subscribe))
(def >evt re-frame.core/dispatch)

(defn obj-name [obj]
  (subs (str/join " " (subvec obj 3)) 1))

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
    (ev/unlisten js/window "mousemove" move-fn)
    (ev/unlisten js/window "mouseup" (:end @state))))

(defn drag-start [id x y evt]
  (let [state   (atom {:pos (client-xy evt)})
        move-fn (drag-move-fn id state)
        done-fn (drag-end-fn id move-fn state)]
    (swap! state assoc :end done-fn)
    (ev/listen js/window "mousemove" move-fn)
    (ev/listen js/window "mouseup" done-fn)
    (>evt [:select-gadget id])
    (.stopPropagation evt)))

(defn get-dom-width [elt]
  (.. (reagent.core/dom-node elt) getBBox -width))

(defn spread-xy [n x y w]
  (let [s (/ (- w 5) (dec n))]
    (mapv #(do [% y]) (range (+ x 2.5) (+ x w) s))))

(defn num-iolets [cmd] ; TODO hard-coded for now
  (case (first cmd)
    :inlet  [0 1]
    :r      [0 1]
    :moses  [2 2]
    :swap   [2 2]
    :metro  [2 1]
    :print  [1 0]
    :outlet [1 0]
    :s      [1 0]
            [1 1]))

(defn iolets-xy [cmd x y w h]
  (let [[ni no] (num-iolets cmd)]
    [(spread-xy ni x y w)
     (spread-xy no x (+ y h) w)]))

(defn obj-as-svg [id x y w h label]
  [:g
    [:rect.obj {:x x :y y :width w :height h}]
    ; https://stackoverflow.com/questions/27602592/reagent-component-did-mount
    ; inlined, this adjusts the enclosing rect's width to the text after render
    [^{:component-did-mount #(>evt [:set-label-width
                                    id
                                    (int (get-dom-width %))])}
      #(do [:text.obj {:x (+ x 5) :y (+ y 15)} label])]])

(defn bang-as-svg [id x y]
  [:circle.bang {:cx (+ x 2.5) :cy (+ y 10) :r 10}])

(defn gadget-as-svg [id [x y typ & cmd :as obj]]
  (let [w          (<sub [:rect-width id])
        h          20
        [ins outs] (iolets-xy cmd x y w h)]
    ^{:key id}
     [:g.draggable {:on-mouse-down #(drag-start id x y %)}
      (map (fn [[cx cy :as xy]]
             ^{:key xy}
              [:circle {:cx cx :cy cy :r 3}]) (concat ins outs))
      (if (= typ :obj)
          (obj-as-svg id x y w h (obj-name obj))
          (bang-as-svg id x y))]))

(defn wire-path [x1 y1 x2 y2] ;; either straight line or cubic bezier
; (str/join " " ["M" x1 y1 "L" x2 y2])
  (str/join " " ["M" x1 y1 "C" x1 (+ y1 25) x2 (- y2 25) x2 y2]))

(defn wire-as-svg [[src-pos src-out dst-pos dst-in :as wire]]
  (let [[sx sy]   (<sub [:gadget-num src-pos])
        [dx dy]   (<sub [:gadget-num dst-pos])
        src-width (<sub [:rect-width src-pos])
        dst-width (<sub [:rect-width dst-pos])]
    ^{:key wire}
     [:path.wire {:d (wire-path (+ sx (* (- src-width 5) src-out) 2.5)
                                (+ sy 20)
                                (+ dx (* (- dst-width 5) dst-in) 2.5)
                                (+ dy 0))}])) 

(defn design-as-svg []
  (let [objs   (<sub [:gadgets])
        wires  (<sub [:wires])]
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
