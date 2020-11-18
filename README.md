# awsbatch
Submit a aws batch job in console command

### How to use

```
Usage:
  awsbatch.exe [OPTIONS]

Application Options:
  /v, /version          show version
      /job-name:        the name of the job
      /job-queue:       the job queue into which the job is submitted
      /job-definition:  the job definition used by this job
      /parameters:      additional parameters.
      /cli-input-json:  performs service operation based on the JSON
      /wait             wait until job finished
      /job-id:          the id of the job

Help Options:
  /?                    Show this help message
  /h, /help             Show this help message
```

### command
Please check bt400-prod.cmd.
<br>
bt400 is a batch to dump postgresql database to file.
