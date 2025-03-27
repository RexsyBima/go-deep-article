package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/alecthomas/kong"
	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/joho/godotenv"
	"strings"
	// token "github.com/pandodao/tokenizer-go"
	"html"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var env = godotenv.Load()

const baseDirDownload = "captions"

type DownloadCmd struct {
	VideoUrls []string `arg:"" required:"" name:"video-url-or-id" help:"Download captions for a YouTube video."`
	Lang      string
}

type Text struct {
	Start   float64 `xml:"start,attr"`
	Dur     float64 `xml:"dur,attr"`
	Content string  `xml:",chardata"`
}

type Transcript struct {
	Texts []Text `xml:"text"`
}

func (d *DownloadCmd) Run() error {
	err := createFolderIfNotExists(baseDirDownload)
	if err != nil {
		log.Fatalf("Unable to create folder: %v\n", err)
	}

	language := d.Lang
	if language == "" {
		log.Println("No language specified, defaulting to English (en)")
		language = "en"
	}

	var wg sync.WaitGroup

	for _, videoUrl := range d.VideoUrls {
		downloader, err := getDownloaderFromArg(videoUrl)
		if err != nil {
			log.Printf("Unable to download captions for %s: %v\n", videoUrl, err)
			continue
		}

		wg.Add(1)
		go downloader.Download(&wg, downloader.DestinationPath(baseDirDownload), language)
	}

	wg.Wait()

	return nil
}

var CLI struct {
	Download DownloadCmd `cmd:"" help:"Download captions for a YouTube video."`
}

func main() {
	args := os.Args
	ctx := kong.Parse(&CLI)
	language := os.Args[len(os.Args)-1]
	toPanic := 0
	for i := 0; i < len(os.Args)-1; i++ {
		if strings.Contains(os.Args[i], "youtube") {
			toPanic += 1
		}
		if toPanic > 1 {
			panic("Program crash, this program only allowed one yt url at a time")
		}

	}

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
	yt_url, err := getParseGetArg(args, "youtube")
	if err != nil {
		panic("Program crash")
	}
	var systemprompt string
	_, isLang, err := extractYTIDFromURL(language)
	fmt.Print(isLang)
	switch isLang {
	case true:
		systemprompt = fmt.Sprintf(`You are an AI transformation agent tasked with converting raw YouTube caption texts about knowledge into a polished, engaging, and readable blog post. Your responsibilities include:
- **Paraphrasing**: Transform the original caption text into fresh, original content while preserving the key information and insights.
- **Structure**: Organize the content into a well-defined structure featuring a captivating introduction, clearly delineated subheadings in the body, and a strong conclusion.
- **Engagement**: Ensure the blog post is outstanding by using a professional yet conversational tone, creating smooth transitions, and emphasizing clarity and readability.
- **Retention of Key Elements**: Maintain all essential elements and core ideas from the original text, while enhancing the narrative to captivate the reader.
- **Adaptation**: Simplify technical details if necessary, ensuring that the transformed content is accessible to a broad audience without losing depth or accuracy.
- **Quality**: Aim for a high-quality article that is both informative and engaging, ready for publication.

Follow these guidelines to generate a comprehensive, coherent, and outstanding blog post from the provided YouTube captions text.
Your final output should be **only** the paraphrased text, styled in Markdown format, and in **%v** language.

`, language)
	default:
		systemprompt = `You are an AI transformation agent tasked with converting raw YouTube caption texts about knowledge into a polished, engaging, and readable blog post. Your responsibilities include:

- Paraphrasing: Transform the original caption text into fresh, original content while preserving the key information and insights.
Structure: Organize the content into a well-defined structure featuring a captivating introduction, clearly delineated subheadings in the body, and a strong conclusion.
- Engagement: Ensure the blog post is outstanding by using a professional yet conversational tone, creating smooth transitions, and emphasizing clarity and readability.
- Retention of Key Elements: Maintain all essential elements and core ideas from the original text, while enhancing the narrative to captivate the reader.
- Adaptation: Simplify technical details if necessary, ensuring that the transformed content is accessible to a broad audience without losing depth or accuracy.
- Quality: Aim for a high-quality article that is both informative and engaging, ready for publication.

Follow these guidelines to generate a comprehensive, coherent, and outstanding blog post from the provided YouTube captions text.
your final output should be only the paraphrased text and style it to markdown like format in english language
	`
	}
	yt_id, _, err := extractYTIDFromURL(yt_url)
	if err != nil {
		panic("Program crash")
	}

	filename := filepath.Join(baseDirDownload, fmt.Sprintf("captions_%s.xml", yt_id))
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		panic("Program crash")
	}
	var transcript Transcript
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&transcript)
	if err != nil {
		panic("Program crash")
	}
	var fulltext string
	for _, text := range transcript.Texts {
		fulltext += html.UnescapeString(text.Content) + " "
	}
	fmt.Printf("data video with id %v has loaded \n", yt_id)
	fmt.Println("generating modified text")
	// tokenCount := token.MustCalToken(fulltext)
	// fmt.Println("Token count: ", tokenCount)

	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: systemprompt},
			{Role: deepseek.ChatMessageRoleUser, Content: fulltext},
		},
	}
	// Send the request and handle the response
	deepseek_ctx := context.Background()
	response, err := client.CreateChatCompletion(deepseek_ctx, request)
	if err != nil {
		panic(err)
	}
	// Print the response
	output := response.Choices[0].Message.Content
	fmt.Println("Response:", output)
	// err = os.WriteFile(fmt.Sprintf("%s/%s.md", baseDirDownload, yt_id), []byte(fulltext), 0644)
	// if err != nil {
	// 	panic(err)
	// }
	filename = fmt.Sprintf("%s/%s.md", baseDirDownload, yt_id)
	file, err = os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(string(output))
	if err != nil {
		panic(err)
	}
	fmt.Printf("finished, the output md file is saved in captions/%v.md\n", yt_id)
	targetFile := fmt.Sprintf("captions_%s.xml", yt_id)
	path := filepath.Join("captions", targetFile)
	err = os.Remove(path)
	if err != nil {
		fmt.Printf("fail to delete redundant filefile", yt_id)
	}
}
