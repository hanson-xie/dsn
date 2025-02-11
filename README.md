<p align="center">
  <a href="https://www.bedrock.technology/" title="bedrock technology">
    <img src="documentation/images/rockx.png" alt="rockx Logo" width="122" />
  </a>
</p>

<h1 align="center">dsn</h1>

## Building & Documentation

> Note: The default `main` branch is the latest stable branch, checkout the most recent [`Latest release`].

## Basic Build Instructions
#### Go

To build dsn client, you need a working installation of [Go 1.23 or higher](https://golang.org/dl/):

```bash
wget -c https://golang.org/dl/go1.23.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```
**TIP:**
You'll need to add `/usr/local/go/bin` to your path. For most Linux distributions you can run something like:

```shell
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc && source ~/.bashrc
```

### Build dsn client

Once all the dependencies are installed, you can build the dsn-cli

1. Clone the repository:

   ```sh
   git clone https://github.com/Bedrock-Technology/dsn.git
   cd dsn/
   ```

Note: The default branch `main` is the pro branch where the latest new features, bug fixes and improvement are in.

2. Build:

   ```sh
   make dsn-cli
   ```
3. Use
    ```sh
   ./dsn-cli --config ./yaml/dsn.yaml run
   ./dsn-cli --config ./yaml/dsn.yaml sql --rpc http://127.0.0.1:8570/dsn/load --toml-file ./toml/redeem_eth.toml
   ```