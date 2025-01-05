## Gitsense

Gitsense is a simple tool designed for GitHub users who want to automate custom actions
triggered by webhooks. While it’s not a full-fledged deployment tool, Gitsense provides
a straightforward way to execute scripts remotely over SSH whenever certain events occur in your
GitHub repositories.

### Getting Started

#### Prerequisite

- A GitHub account with access to the repositories you want to manage.
- A Github ```oauth2``` app.
- A server or remote machine accessible via SSH where your scripts will be executed.
- The following information for your remote machine:
    - IP address
    - Username
    - Private key
    - Host key
- An environment ready for running Gitsense (e.g., a local or cloud server hosting Gitsense), with Go installed.

### Setting Up Gitsense

#### Clone the repository

```shell
git clone https://github.com/yourusername/gitsense.git
```

#### Navigate to the directory

```shell
cd gitsense
```

#### Deploying Gitsense

The respository has a custom build and deployment script which uses ```systemd``` and ```Caddy```.
If you are familiar with these tools, consider using it. Create a Caddyfile:

```shell
sudo nano /etc/caddy/Caddyfile
```

```
<your_domain_name> {
    handle {
        reverse_proxy localhost:3000
    }
}
```

Or, if you are also planning to use the ```gitsense-client```, a web app built with React.
It has a simple interface that uses ```gitsense```, use:

```
<your_domain_name> {
    handle /api/* {
        reverse_proxy localhost: 3000 // gitsense
    }

    handle {
        reverse_proxy localhost: 4173 // gitsense-client
    }
}
```

Both ```gitsense``` and ```gitsesnse-client``` has a ```./build.sh``` script that
runs them in the background as a service.

#### Make the build script executable:

```shell
chmod +x ./build.sh
```

#### Run the build script:

```shell
./build.sh
```

#### Database

Set-up a ```mongodb``` atlas and add the connection URI in the ```.env``` file:

```shell
MONGODB_URI=<connection_uri>
```

#### Environment Variables

##### Gitsense

- ```PORT``` Port where  the server will run

- ```VERSION``` Version that will be injected to the API endpoints, set it to ```v1```

- ```GITHUB_CLIENT_ID``` Used to identify your Github ```oauth2``` app

- ```GITHUB_CLIENT_SECRET```

- ```GITHUB_RERDIRECT_URL``` Github ```oauth2```'s callback url, set it to ```https://your_domain.com/api/v1/auth/github/callback```

- ```GITHUB_API_URL``` Github's base api URL, ```https://api.github.com```

- ```MONGODB_URI``` Your mongoDB atlas' connection URI

- ```DOMAIN``` Used for setting the cookies, set it to your domain, ```your_domain.com```

- ```BASE_SERVER_URL``` Your server's URL, ```https://your_domain.com```, used to construct the webhook's endpoint

- ```BASE_CLIENT_URL``` Your client's URL, ```https://your_domain.com```, used for redirect when authenticated

- ```ENCRYPTION_KEY``` Used by the ```aes``` algorithm for encryption and decryption

     Both ```refresh``` and ```access``` token are only used to access data from the database, where webhook's configurations are stored including the ```ssh``` configuration.

- ```REFRESH_TOKEN_KEY``` 

- ```ACCESS_TOKEN_KEY```


##### Gitsense Client

- ```BASE_API_URL``` Set it to ```https://your_domain.com/api/v1```
- ```PORT``` Port where  the server will run, default is ```4173```

### Creating a Webhook

Once Gitsense is up and running, you can create webhooks for your GitHub repositories.

- Log in at `gitsense-client` using your Github account

- Select a repository, each repository has a custom context menu; right click and select add webhook

#### Configuring a ```Webhook```

- ```Name``` A unique name that identifies your ```webhook```

- ```Secret``` Used to verify the request coming in for the ```webhook```

- ```IP address``` IP address of  the remote machine
    - Used to identify your remote machine
    - to get your remote machine's IP open the terminal and execute:
    
        ```shell
        curl ifconfig.me
        ```
- ```Private key``` The private key used to ```ssh``` into your remote machine
    - Typically stored in ```.ssh``` directory, in your local machine open the terminal
    and execute:

        ```shell
        cd .ssh
        ```

        ```shell
        ls
        ```
    - The private key is encrypted using ```aes``` algorithm.

- ```Host key``` Host key of the remote machine
    - Host key will be used to verify the ssh connection, and prevent ```MITM``` attacks.
    
    - Typically this is stored on the remote machine's ```/etc/ssh``` directory, open the terminal and execute:

        ```shell
        cd /etc/ssh
        ```

        ```shell
        ls
        ```
        Find out which algorithm is used on your private key and select the approriate one.

        ```shell
        sudo cat ssh_host_<algorithm>_key.pub
        ```

- ```Script``` The script that will be executed on the remote machine
    - Sample script:

        ```shell
        cd <project_directory>
        git pull
        ./build.sh
        ./deploy.sh
        ```

### Development Mode

To run on development mode, make sure you got ```air``` installed, on gitsense's root directory execute:

```shell
air
```

#### Test an SSH Connection

##### SSH Configuration

To test an ```ssh``` connection, set your ssh configuration in the ```.env``` file:

```
PRIVATE_KEY=<private_key>
HOST_KEY=<host_key>
INSTANCE_IP=<ip_address>
USER=<username>
```

##### Change directory

```shell
cd pkg/util/ssh
```

##### Run test

```shell
go test
```












