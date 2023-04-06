package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/louisbuchbinder/aime"
)

var (
	temperature     float64
	model           string
	prompt          string
	systemPrompt    string
	systemPromptKey string
	key             string
	keyfile         string
	output          string
	echo            bool
)

func init() {
	flag.BoolVar(&echo, "echo", false, "Echo the prompt")
	flag.Float64Var(&temperature, "temperature", 0, "Temperature affects the randomness of the response. Range 0 <= t <= 2, where 2 is most random. The default is 0.")
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "The OpenAI model to use. The default is 'gpt-3.5-turbo'")
	flag.StringVar(&prompt, "prompt", "", "The literal prompt. If empty then STDIN data will be used")
	flag.StringVar(&systemPrompt, "systemPrompt", "", "The literal systemPrompt. If empty then the systemPromptKey will be used to determine the system prompt")
	flag.StringVar(&systemPromptKey, "systemPromptKey", "", "The system prompt flag. If empty then no system prompt will be used. If provided then a system prompt will be determined for the key")
	flag.StringVar(&key, "key", "", "The OpenAI API key. Oneof key or keyfile must be provided")
	flag.StringVar(&keyfile, "keyfile", "", "The path to a file containing the OpenAI API key. Oneof key or keyfile must be provided")
	flag.StringVar(&output, "output", "raw", "The output format. Supports 'raw', 'inline'. Defaults to 'raw'")
	flag.Parse()

	if key == "" && keyfile == "" {
		panic(fmt.Errorf("Oneof key or keyfile must be provided"))
	}
	if output != "raw" && output != "inline" {
		panic(fmt.Errorf("output must be 'raw' or 'inline', instead got: %s", output))
	}
}

func main() {
	if prompt == "" {
		if b, err := ioutil.ReadAll(os.Stdin); err != nil {
			panic(err)
		} else {
			prompt = string(b)
		}
	}

	if systemPrompt == "" {
		systemPrompt = aime.LookupSystemPrompt(systemPromptKey)
	}

	var messages = []*aime.Message{}

	if systemPrompt != "" {
		messages = append(messages, &aime.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	messages = append(messages, &aime.Message{
		Role:    "user",
		Content: prompt,
	})

	req, err := aime.ToRequest(&aime.RequestData{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	})
	if err != nil {
		panic(err)
	}

	var opts = []aime.ClientOption{}
	if key != "" {
		opts = append(opts, aime.WithKey(key))
	} else {
		opts = append(opts, aime.WithKeyFile(keyfile))
	}

	client, err := aime.ToClient(opts...)
	if err != nil {
		panic(err)
	}

	data, err := aime.MakeRequest(client, req)
	if err != nil {
		panic(err)
	}

	if echo {
		fmt.Printf("%s\n%s\n\n", systemPrompt, prompt)
	}

	if output == "inline" {
		for _, c := range data.Choices {
			fmt.Println(c.Message.Content)
		}
	} else {
		if b, err := json.Marshal(data); err != nil {
			panic(err)
		} else {
			fmt.Println(string(b))
		}
	}

}
