package parsers

import "testing"

func TestGroupedSimilarClosedEpisodesAndPack(t *testing.T) {
	page := `<h1 class="meta__title">Example</h1><details><summary>Season 2</summary><details><summary>Episode 9</summary><table><tbody>
 <tr><th class="similar-torrents__type" rowspan="2">WEB-DL</th><td><h3 class="torrent-search--grouped__name"><a href="https://utp.to/torrents/101">2160p-GROUP</a></h3></td><td class="torrent-search--grouped__size"><span title="9122854191 B">8.5 GiB</span></td><td class="torrent-search--grouped__seeders">51</td></tr>
 <tr><td><h3 class="torrent-search--grouped__name"><a href="https://utp.to/torrents/102">1080p-OTHER</a></h3></td></tr>
 </tbody></table></details></details>
 <details><summary>Season 1</summary><table><tbody><tr><td><h3 class="torrent-search--grouped__name"><a href="/torrents/103">1080p pack</a></h3></td></tr></tbody></table></details>`
	got, err := ParseSimilar(0, page)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || len(*got) != 3 {
		t.Fatalf("similar = %#v", got)
	}
	if (*got)[0].ID != 101 || (*got)[0].SizeBytes != 9122854191 || (*got)[0].Seeders != 51 || (*got)[1].Type != "WEB-DL" {
		t.Errorf("rows = %#v", *got)
	}
}

func TestSimilarRejectsLogin(t *testing.T) {
	if _, err := ParseSimilar(0, `<form action="/login"></form>`); err == nil {
		t.Error("accepted login")
	}
}
