package controllers

import (
	"github.com/revel/revel"
	"www-forecast/app/models"
	"regexp"
	"sort"

	)

type RankResults struct {
	*revel.Controller
}

type Votes []Vote

type Vote struct {
	Index int
	Votes int
}

func (v Votes) Len() int		{ return len(v) }
func (v Votes) Swap(i, j int)		{ v[i], v[j] = v[j], v[i] }
func (v Votes) Less(i, j int) bool { return v[i].Votes < v[j].Votes}


func (c RankResults) Index(rid string) revel.Result {

		c.Validation.Required(rid)
		c.Validation.Match(rid, regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"))
		//^[a-e0-9]{8}-[a-e0-9]{4}-[a-e0-9]{4}-[a-e0-9]{12}$

		if c.Validation.HasErrors() {
			c.Flash.Error("Cannot view results, errors in submission.")

			return c.Redirect(List.Index)
		}

		//This attempts to retrieve the scenario based on the hosted domain, for security.
		r := models.GetRank(rid)

		// We use the SID from the successful call using the hosted domain, instead of whatever the user gives us.
		rr := models.ViewRankResults(r.Rid)
		if (len(rr)>0){
			pw := getPositionalWinner(rr)
			//Return rank results, rank, copeland rank
			return c.Render(rr, r, pw)
		} else {
			c.Flash.Error("No results yet.")
			return c.Redirect("/view/estimate/%s", rid)
		}

}

func getPositionalWinner(rr []models.Sort) (vs Votes){

	vs = make(Votes, len(rr[0].Options))

	total := len(rr[0].Options)

	//First loop. 'v' is a Sort.
	for _, v := range rr {

		//Second loop. Each "o" is a preference, top to bottom.
		for i, o := range v.Options {
			vs[o].Votes += total - i - 1
			vs[o].Index = o
		}
	}

	sort.Sort(sort.Reverse(vs))

	return vs
}
