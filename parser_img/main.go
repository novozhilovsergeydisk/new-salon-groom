package main

import (
 "fmt"
 "net/http"
 "os"

 "golang.org/x/net/html"
)

// findImages находит все изображения в HTML-документе
func findImages(n *html.Node) []string {
 var images []string
 if n.Type == html.ElementNode && n.Data == "img" {
  for _, attr := range n.Attr {
   if attr.Key == "src" {
    images = append(images, attr.Val)
   }
  }
 }
 for c := n.FirstChild; c != nil; c = c.NextSibling {
  images = append(images, findImages(c)...)
 }
 return images
}

// getHTML загружает HTML-код страницы по указанному URL
func getHTML(url string) (*html.Node, error) {
 resp, err := http.Get(url)
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
  return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
 }

 doc, err := html.Parse(resp.Body)
 if err != nil {
  return nil, err
 }
 return doc, nil
}

func main() {
 url := "https://top-tomcat.10web.me/"
 doc, err := getHTML(url)
 if err != nil {
  fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
  return
 }

 images := findImages(doc)
 fmt.Println("Found images:")
 for _, img := range images {
  fmt.Println(img)
 }
}
