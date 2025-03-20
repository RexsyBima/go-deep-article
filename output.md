# Go Deep Article

A program that converts YouTube video content into paraphrased articles, written in Go (heck yea). Perfect for content branding in texted media (medium, linkedin and others). Transforms spoken content into written format while maintaining core meaning.


## Installation and uses
first, make sure your env system has api key for deepseek, named `DEEPSEEK_API_KEY`
Install go-deep-article by cloning my repo and run it with go

```bash
git clone https://github.com/RexsyBima/go-deep-article.git
cd go-deep-article
go run . download "youtubevideourlhere"
```

  OR

download the binary file in the release here, then run `deeparticle download "youryoutubeurlhere"`

your paraphrased md file will be named by your youtube video id in folder `captions` in the same level of where you run the program from 
## Features

- Converts YouTube video transcripts into readable articles
- Paraphrases content while preserving key information
- Output into markdown file
## ⚠️Important Disclaimers

**User Responsibility Notice**:  
This tool requires **human verification** of all outputs. Users **must**:
- Review and understand generated content before use
- Ensure the transformed article respects copyright and fair use principles
- Verify factual accuracy of the information
- Never present AI-generated content as human-written without disclosure

**Liability Disclaimer**:  
THE DEVELOPER MAKES NO WARRANTIES ABOUT THE OUTPUT QUALITY OR ACCURACY. BY USING THIS TOOL, YOU AGREE THAT:
- You bear sole responsibility for any content created with this tool
- The developer is not liable for any misuse, legal issues, or damages arising from generated content
- Outputs may contain errors or inaccuracies that require human correction
