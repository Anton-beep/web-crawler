package analyser

const KeyWordsPromptPrefix = "Extract the main keywords from the following texts and provide a list of keywords SEPARATED BY NEW LINES(\\n). The response must contain only the keywords as a list, nothing else. Answer in English, regardless of the language of the input texts.\n\nInput texts:\n"
const KeyWordsPromptSuffix = "\nExtracted Keywords:\n"

const MainIdeaPromptPrefix = "Extract the main idea from the following texts and provide a single paragraph summarizing the key information. The response must consist of only one paragraph with no subheadings or additional formatting. Answer in English, regardless of the language of the input texts.\n\nInput texts:\n"
const MainIdeaPromptSuffix = "\nMain Ideas:\n"
