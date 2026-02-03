## Input

[./input.txt](./input)

Attempted to use a handful of techniques but the input file is a combination of #1 and #2. Once subdomains were found they were validated by making get requests to {subdomain}/careers and {subdomain}/jobs, which seemed to be the most common endpoints. If either of those returned a 200 response, the response url (because some of these redirects to the correct careers page) was then added to the final input.

1. Fetched indexes from commoncrawl.org querying each for urls that matched \*.avature.net
2. Used a pentest tool to query subdomains for avature.net
3. Queried crt.sh for matching domains of \*.avature.net. This retuened lackluster results.
4. Queried search engines for subdomains. Search results are hidden with javascript which we would need a headless browser to get around.

## Output

[output.txt](https://drive.google.com/file/d/1G0_wHPs5pY7aBJp6XjFcvhPeDs1Z2AkI/view?usp=sharing)

Jobs were mostly scraped using the most commonly seen css classes. The flow was to navigate to the jobs page, grab the next link, grab all jobs on the page scraping the title and details link, navigate to each detail page to get the job details and metadata, and then continue on to the next page if there was one.

## Improvements

- Concurrency - scraping took a few hours so adding the ability to fetch multiple sites at once would be beneficial.
- Add headless browser - some sites and jobs are missed because they simply require javascript to prevent scraping so using a headless browser would circumvent this.
