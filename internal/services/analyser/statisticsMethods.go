package analyser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Response struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	TotalDuration int    `json:"total_duration"`
	LoadDuration  int    `json:"load_duration"`
}

func GetAllAnalyseMethods() []string {
	return []string{"KeyWords", "MainIdeas"}
}

func (s *Server) sendToLLM(model string, prompt string) (Response, error) {
	request := Request{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	body, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", s.Analyser.Link, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	ctx := context.Background()
	req = req.WithContext(ctx)
	client := &http.Client{}
	ans, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	bit, err := io.ReadAll(ans.Body)
	if err != nil {
		return Response{}, err
	}
	resp := Response{}
	err = json.Unmarshal(bit, &resp)
	if err != nil {
		return Response{}, err
	}
	return resp, nil
}

func createDelimiter(texts []string) string {
	toReturn := ""
	for n, text := range texts {
		toReturn += fmt.Sprintf("TEXT %v\n%s", n+1, text)
	}
	return toReturn
}

func (s *Server) keyWordsPrompt(texts []string) (Response, error) {
	resp, err := s.sendToLLM("mistral-nemo:latest", Role+KeyWorldPromptFormat+Delimiter+Audition+"\n"+createDelimiter(texts))
	if err != nil {
		return Response{}, err
	}
	return resp, nil
}

func (s *Server) MainIdeasPrompt(texts []string) (Response, error) {
	resp, err := s.sendToLLM("mistral-nemo:latest", Role+MainIdeasPromptFormat+Delimiter+Audition+"\n"+createDelimiter(texts))
	if err != nil {
		return Response{}, err
	}
	return resp, nil
}
