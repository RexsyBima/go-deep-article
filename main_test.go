package main_test

import (
	"context"
	"encoding/xml"
	"fmt"
	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/joho/godotenv"
	token "github.com/pandodao/tokenizer-go"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var env = godotenv.Load()

type Text struct {
	Start   float64 `xml:"start,attr"`
	Dur     float64 `xml:"dur,attr"`
	Content string  `xml:",chardata"`
}

type Transcript struct {
	Texts []Text `xml:"text"`
}

func TestArgParsing(t *testing.T) {
	t.Skip()
	expected := "indo"
	// simulate os args
	var args []string
	args = append(args, "value1", "value2", "value3", "https://www.youtube.com/watch?v=Fo49GokDJhM", "indo")
	for _, arg := range args {
		if strings.Contains(arg, expected) {
			fmt.Printf("argumen found")
			return
		}
	}
	t.Errorf("fail to parse expected value")
}

func TestRegexYTIDExtract(t *testing.T) {
	t.Skip()
	url := "https://www.youtube.com/watch?v=Fo49GokDJhM&t=25"
	re := regexp.MustCompile(`(?:v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		fmt.Println("Video ID:", match[1])
	} else {
		t.Errorf("fail to parse regex expression value")
	}
}

func TestPathJoin(t *testing.T) {
	t.Skip()
	dir := "captions"
	vid_id := "captions_IbXSEGB8LRs.xml"
	output := filepath.Join(dir, vid_id)
	if output != "captions/captions_IbXSEGB8LRs.xml" {
		t.Errorf("fail to join path")
	}
}

func TestXMLParsing(t *testing.T) {
	t.Skip()
	filename := filepath.Join("captions", "captions_IbXSEGB8LRs.xml")
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		t.Errorf("fail to open file")
	}
	var transcript Transcript
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&transcript)
	if err != nil {
		t.Errorf("fail to parse xml file")
	}
	var fulltext string
	for _, text := range transcript.Texts {
		fulltext += html.UnescapeString(text.Content) + " "
	}

}
func TestTokenizer(t *testing.T) {
	t.Skip()
	mdInput := "captions/sVC4DL2secQ.md"
	data, err := os.ReadFile(mdInput)
	if err != nil {
		t.Errorf("fail to read file: %v", err)
	}
	tokenCount := token.MustCalToken(string(data))
	fmt.Println(tokenCount)
}

func TestDeepseekGo(t *testing.T) {
	t.Skip()
	filename := filepath.Join("captions", "captions_IbXSEGB8LRs.xml")
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		t.Errorf("fail to open file")
	}
	var transcript Transcript
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&transcript)
	if err != nil {
		t.Errorf("fail to parse xml file")
	}
	var fulltext string
	for _, text := range transcript.Texts {
		fulltext += html.UnescapeString(text.Content) + " "
	}
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	// Create a chat completion request
	systemprompt := `You are an AI transformation agent tasked with converting raw YouTube caption texts about knowledge into a polished, engaging, and readable blog post. Your responsibilities include:

- Paraphrasing: Transform the original caption text into fresh, original content while preserving the key information and insights.
Structure: Organize the content into a well-defined structure featuring a captivating introduction, clearly delineated subheadings in the body, and a strong conclusion.
- Engagement: Ensure the blog post is outstanding by using a professional yet conversational tone, creating smooth transitions, and emphasizing clarity and readability.
- Retention of Key Elements: Maintain all essential elements and core ideas from the original text, while enhancing the narrative to captivate the reader.
- Adaptation: Simplify technical details if necessary, ensuring that the transformed content is accessible to a broad audience without losing depth or accuracy.
- Quality: Aim for a high-quality article that is both informative and engaging, ready for publication.
Follow these guidelines to generate a comprehensive, coherent, and outstanding blog post from the provided YouTube captions text.

	your final output should be only the paraphrased text and style it to markdown like format
	`
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: systemprompt},
			{Role: deepseek.ChatMessageRoleUser, Content: fulltext},
		},
	}
	// Send the request and handle the response
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		t.Errorf("fail to create chat completion: %v", err)
	}
	// Print the response
	output := response.Choices[0].Message.Content
	fmt.Println("Response:", output)
	err = os.WriteFile("output.md", []byte(fulltext), 0644)
	if err != nil {
		panic(err)
	}
}

func TestSaveMdData(t *testing.T) {
	t.Skip()
	mdFileName := "README.md"
	mdOutputTestFile := "output.md"
	data, err := os.ReadFile(mdFileName)
	if err != nil {
		t.Errorf("fail to read file: %v", err)
	}
	fmt.Println(string(data))
	file, err := os.Create(mdOutputTestFile)
	if err != nil {
		t.Errorf("fail to create file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(string(data))
	if err != nil {
		t.Errorf("fail to write file: %v", err)
	}
}

func TestBooleanLogic(t *testing.T) {
	t.Skip()
	var val1 string
	fmt.Print(val1 == "")
}

func TestReverseParsingArg(t *testing.T) {
	var args []string
	expected := "indo"
	args = append(args, "value1", "value2", "value3", "https://www.youtube.com/watch?v=Fo49GokDJhM", "indo")
	final_arg := args[len(args)-1]
	re := regexp.MustCompile(`(?:v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	match := re.FindStringSubmatch(final_arg)
	if len(match) > 1 {
		t.Errorf("fail to parse regex expression value")
	}
	fmt.Println("language id:", final_arg == expected)
}

func TestDeleteFile(t *testing.T) {
	t.Skip()
	targetFile := "captions_A3IQj0LCocA.xml"
	path := filepath.Join("captions", targetFile)
	err := os.Remove(path)
	if err != nil {
		t.Errorf("fail to delete file: %v", err)
	}
}

func TestPointer(t *testing.T) {
	data := "foobar"
	p * string := &data
}
