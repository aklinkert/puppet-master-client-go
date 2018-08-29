# puppet-master-client-go

Golang SDK for the [puppet-master.io](https://puppet-master.io) public API. Puppet-master makes the execution of website interactions
super simple by abstracting the code execution behind a HTTP API, scheduling the job for you in a scalable
manner. For more information please head over to the [puppet-master docs](https://docs.puppet-master.io).


## installation

```bash
go get github.com/Scalify/puppet-master-client-go
```

## example usage

````go
package main

import (
	"log"

	"github.com/scalify/puppet-master-client-go"
)

func main() {
	client, err := puppetmaster.NewClient("https://puppet-master.io/api/v1/teams/my-team", "theapitokenigot")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	client.EnableDebugLogs()

	jobs, err := client.GetJobs(1, 100)
	if err != nil {
		log.Fatalf("failed to fetch jobs: %v", err)
	}

	log.Printf("Current page %d, last page %d, total %d", jobs.Meta.CurrentPage, jobs.Meta.LastPage, jobs.Meta.Total)

	for _, job := range jobs.Jobs {
		log.Printf("Job ID %v", job.UUID)
	}

	newJob := &puppetmaster.JobRequest{
		Code: `
import {getIp} from 'shared';

await page.goto(vars.page);
const ip = await getIp(page);

logger.info(ip);
results.ip = ip;
`,
		Modules: map[string]string{
			"shared": `
export async function getIp(page) {
  const text = await page.evaluate(() => document.querySelector('body').textContent);
  return text.split(":")[1];
}
`,
		},
		Vars: map[string]string{
			"page": "http://ifcfg.co",
		},
	}

	createdJob, err := client.CreateJob(newJob)
	if err != nil {
		log.Fatalf("failed to create job: %v", err)
	}

	log.Printf("Created job %v", createdJob.UUID)

	retrievedJob, err := client.GetJob(createdJob.UUID)
	if err != nil {
		log.Fatalf("failed to retrieve single job: %v", err)
	}

	log.Printf("Retrieved job: %v", retrievedJob.UUID)

	err = client.DeleteJob(retrievedJob.UUID)
	if err != nil {
		log.Fatalf("failed to delete job: %v", err)
	}

	log.Printf("Done executing lifecycle of job %v", retrievedJob.UUID)
}
````

## License

Copyright 2018 Scalify GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
