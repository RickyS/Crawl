Program to crawl the web. Uses companion package creep.
=======================================================

Given a website or list of websites, fetch the web page, and then fetch all the pages linked to
by that website.  
And so on?  Sometimes.  See below.

To install:  
       $ go get github.com/RickyS/crawl  
       $ go get github.com/RickyS/creep  

You'll neeed both packages, the depend on each other.  The main program is crawl. 
The working package is creep.

To run:  
      $ go run crawl.go  
Or, on Linux:  
      $ go install -v -x  && crawl  

Instead of command line arguments, the crawl program reads a json file.
A simple small file is the default:  
iana.json, which is among the included files.

crawl.go contains the line:  
     jobData := creep.LoadJobData("**iana.json**")  
which specifies that the parameter file iana.json will be used.  

Todo: The name of the json file should be a command line argument.

Technically, crawling the web isn't possible.  The number of websites is effectively infinite.  So after
numerous tests that filled up my machine, I have artificially truncated the number of web pages crawled.

There are parameters in the json file that adjust the limitations.  TBD.
