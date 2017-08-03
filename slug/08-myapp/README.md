# 08-myapp

A fairly comfortable development environment, with Reagent (React wrapper) and
Devtools (nice printouts in Chrome's JavaScript console). Project was generated
with `boot -d boot/new new -t tenzing -n myapp -a +reagent -a +devtools`.

To run in dev mode:

    boot dev

Once you see the `Elapsed time: ...` message, open <http://localhost:3000>.

To launch a REPL, run this in the separate terminal window:

    boot repl -c

Then, at the `boot.user` prompt, start the REPL connection to the browser:

    boot.user=> (start-repl)

Some things to try at the REPL:

    cljs.user=> (js/alert "123")
    nil
    cljs.user=> (ns myapp.app)

    myapp.app=> (some-component)
    WARNING: Use of undeclared Var myapp.app/some-component at line 1 <cljs repl>
    [:div [:h3 "I am a component!"] [:p.someclass "I have " [:strong "bold"] [:span {:style {:color "red"}} " and red"] " text."]]
    myapp.app=>

The warning may have something to do with running two JVMs? (still learning...)
