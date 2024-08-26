# UmaAIChatServer

This project is the server part of [UmaAIChat](https://github.com/KLXLjun/UmaAIChat), which lets you chat with Umamusume. Any updates will be reflected in the to-do list under [UmaAIChat](https://github.com/KLXLjun/UmaAIChat). 

> This project supports using large language models compatible with the OpenAI API, like Ollama.

## How to Use It?

1. Get the latest version from the [Releases](https://github.com/KLXLjun/UmaAIChat/releases) page.
2. After running it for the first time, a file named `config.toml` will be created in your program's directory. Close the program once this is done.
3. In the configuration file:

   - Under `chat` and `emotion`, adjust the settings according to your needs.

   `chat` deals with the API setup for chatting with the large language model, while `emotion` handles setting up an API for emotion recognition.

   Here's what each option means:
   
   - `api_url`: The URL of the API.
   - `api_key`: Your API key. (You need to obtain this from a service that provides the compatible models)
   - `model`: Choose the model you want to use.
   - `proxy` (Optional): If you need an HTTP proxy, enter its address here.

4. Restart the program after making your changes.

For chatting, it's best to use either `gpt-4o` or `gpt-4o-mini`. For emotion recognition, use a local model like Ollama's quantized model.

