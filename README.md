## Basic usage

The original aim of this package is to help scrapping expat.cz articles as such:

```golang
func main() {
	link, err := scraper.FindLinkWith("weekend", scraper.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	ct, err := scraper.GetArticleContent(link)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("events.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewEncoder(file).Encode(ct)
	if err != nil {
		log.Fatal(err)
	}
}
```
giving result like this
```json
{
    "Content": "On May 18, peruse military equipment and weapons at Atrium Flora, which will be transformed into a military base...",
    "Title": "MILITARY MARVELS "
}
```
