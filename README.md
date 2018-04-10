# Microb http

Http service for [Microb](https://github.com/synw/microb). Features:

- **Serve html** from the filesystem (other datasources are planned for later)
- **Record hits** into an sqlite database with geoip info
- **Autoreload** pages and templates for development

As this is a Microb service it has remote commands and records logs into an sqlite database

Requirements: 
[Centrifugo](https://github.com/centrifugal/centrifugo/) and [Redis](http://redis.io/)

#### Install and status

To install you have to compile Microb with the http service 
for now as no release has been made. The dev status is move fast and break things for the moment.

To autoreload pages and templates during development start the instance 
with the `-d` flag: `./microb -d`

To start the http server without the Microb client: `./microb -s`

#### Configuration

Configure `config_http.json`:

   ```javascript
   {
	   "centrifugo_addr":"localhost:8001",
	   "centrifugo_key":"secret_key",
	   "domain": "localhost",
	   "addr":"localhost:8080",
	   "mail":false,
	   "csrf_key":"297gcbbb-2b24-4262-b6f7-817e0983fbb7",
	   "hitsDb": "hits.sqlite"
   }
   ```

## Usage

### Templates

The templates are in the `templates` folder. These are editable Go templates.

**Note**: when the templates are edited they will be automatically reparsed, no need for a
server restart

### Content

To serve pages using the filesystem create content in `static/content`. Ex: `index.html`:

   ```html
   <!-- Title:My homepage -->
   <h1>Home</h1>
   ```
   
The `<!-- Title:My homepage -->` will be used to populate the `<title>` head tag.

The urls will automatically map to the filesystem structure: `static/content/folder/page.html` 
will be accessible at `/folder/page` 

### Static files

All that is in the `static` folder will be accessible at `/static`

### Hot reload

If the dev mode is enabled any change in templates, content or static files will
trigger a reload in the browser

## Cli commands

`ping`: check if the http server is up

`status`: tells if the server is running

`start`: start the http server

`stop`: stop the http server
