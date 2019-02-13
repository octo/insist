# insist

**insist** retries commands until they succeed.

I'm frequently travelling with the train where the internet connection is
spotty. Commands like `git push` run into a timeout and I have to restart them,
without any changes to the command line. Computers are much better suited for
menial tasks like this.

**Example:**

```
octo:~/collectd$ ~/go/bin/insist -- git remote update github
2019/02/13 09:03:25 Attempt 0 ...
Fetching github
```

By default, *insist* launches the program and waits for it to finish. If the
program exits with a non-zero exit status, an exponential back-off is applied
and then the program is restarted. You can limit the number of attempts and
provide a timeout via command line options.

If `STDIN` is not a terminal, *insist* copies it into a buffer and passes the
buffer to the program on each attempt. This allows piping data into a retried
program. Run interactive programs with *insist* will not work, because the
program will only be executed once all the input is read.

## License

*insist* is free software. It is licensed under the
[ISC License](https://opensource.org/licenses/ISC).
