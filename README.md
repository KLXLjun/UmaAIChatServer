# UmaAIChatServer

[English document here](https://github.com/KLXLjun/UmaAIChatServer/blob/master/README-en.md)

此项目是 [UmaAIChat](https://github.com/KLXLjun/UmaAIChat) 的服务端，一起使用就可以和马娘进行对话。更新计划会合并在 [UmaAIChat](https://github.com/KLXLjun/UmaAIChat) 下方代办列表中。

> 此项目支持调用兼容于`OpenAI API`的大语言模型（例如`Ollama`等）。

## 如何使用？

1. 从 [Releases](https://github.com/KLXLjun/UmaAIChat/releases) 页面下载最新的程序。

2. 在运行第一次过后，程序目录下会生成一个名为`config.toml`的配置文件，文件生成之后关闭程序。

3. 配置文件中，`chat`和`emotion`这两个大类中的参数根据你的情况进行修改。

   `chat`代表用于聊天的大语言模型API配置，`emotion`是用于识别情绪的大语言模型API配置。

   配置项说明如下：

   - `api_url`: API地址
   - `api_key`：API密钥
   - `model`：使用的模型名称
   - `proxy`：代理地址（留空即为不设置代理）
   
4. 修改完成之后再启动即可使用。

推荐使用`gpt-4o`或是`gpt-4o-mini`作为聊天使用模型，情绪识别则使用本地模型（例如`Ollama`的量化模型之类的）

