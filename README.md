# hugo.drive
A bridge between Google Documents and Hugo

## Overview

The goal of this project was to allow Google Documents to act as content files for Hugo projects. [Hugo](http://gohugo.io/) is an excellent static site generator that allows me to quickly create robust websites that need virtually no CPU overhead to operate. In the news business, we create several new projects a month. Hugo allows me to develop them, get them out the door, and then forget about them so I can move on to the next task.

One of the major drawbacks of Hugo is that it requires files be saved directly on the server. Like many other businesses, we have vendors, management systems, and preferred employee workflows that don't mesh well with a remote repository of code. We heavily rely on Google Drive, however. By building a bridge between Hugo and Drive, I can create the necessary content files automatically using a cron, and our employees can comofortably work on a platform that they know well.

The `hugo.drive` command simply saves Google Documents, by folder, into corresponding folders inside the hugo project. I chose Go to both learn the language and provide a way to run this program on both Mac and Linux. This is my first Go program, and I used quite a bit of code from the [Google documentation](https://developers.google.com/drive/v3/web/quickstart/go).

## Configuration/Flags

**-secret string (default "hd.secret.json")**

To authenticate, you must create a Google App and download a JSON credential file. The [Google documentation](https://developers.google.com/drive/v3/web/quickstart/go) tells you how to do this. Rename that file to *hd.secret.json* and place it in your current working directory, or use the -secret flag to pass a filepath to your favorite place.

When you first run the program, it will provide a url to paste into your browser. After authenticating, the token response will be saved to $HOME/.hugo.drive/auth-token.json. **It's important to know what user authenticated if you want to put this on a cron later**. This is a one-time process and all but one simplification function came directly from Google.


**-config string (default "hd.config.json")**

You must create a configuration file for the program to run. You can either name the file *hd.config.json* and place it in the current working directory, or you can pass in a filepath to wherever you like using the -config flag. Below is an example of the structure.

```json
{
  folders: [
    {
      "id": "*ID of the folder in Drive*",
      "path": "*path of the folder on the local system (e.g. /content/articles)*"
    }
  ]
}
```


**-modified bool**

The modified tag will create the file *hd.modified.txt* in the current working directory containing the RFC 3339 timestamp of the last run. Each additional pull will only look for files in the configured folders that have been modified since that time, and update the *hd.modified.txt* file.

## Final thoughts

I tried to keep this process separate from Hugo itself, partly because you have to but also because it just feels right. I like being able to mix and match what sections are linked to Drive and what sections are static. Sure, it means two commands to build the site each time, but a tiny Makefile works really well for me.

If you don't have Go setup on your machine, I have included binaries for my two targets in the /bin folder.

If you are an avid Go programmer and you see some antipatterns or something I did wrong, please let me know. I'm very new and really interested.
