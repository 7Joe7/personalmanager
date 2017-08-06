package hello

import (
    "net/http"
    "appengine"
    "appengine/datastore"

    "github.com/7joe7/personalmanager/resources"
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
}

// datastoreKey returns the key used for all guestbook entries.
func datastoreKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "PEMA", "default_pema", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    // Ancestor queries, as shown here, are strongly consistent with the High
    // Replication Datastore. Queries that span entity groups are eventually
    // consistent. If we omitted the .Ancestor from this query there would be
    // a slight chance that Greeting that had just been written would not
    // show up in a query.
    q := datastore.NewQuery("Habit").Ancestor(datastoreKey(c)).Limit(10)
    habits := make([]resources.Habit, 0, 10)
    if _, err := q.GetAll(c, &habits); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func sign(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    g := resources.Habit{
        Name: "test",
        Active: true,
    }
    // We set the same parent key on every Greeting entity to ensure each Greeting
    // is in the same entity group. Queries across the single entity group
    // will be consistent. However, the write rate to a single entity group
    // should be limited to ~1/second.
    key := datastore.NewIncompleteKey(c, "Habit", datastoreKey(c))
    _, err := datastore.Put(c, key, &g)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}