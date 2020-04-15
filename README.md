# Get Me My Flour

As I'm aspiring baker/pizza maker, I like to have my flour supplies uninterrupted. Sadly, nowadays, that's not that easy.

## Usage

Prereqs:
* `Go`
* `Gmail` account

* Set up these environmental variables:
  * `FREQUENCY`
  * `NOTIFICATION_EMAIL_SEND`
  * `NOTIFICATION_EMAIL_SEND_PASSWORD`
  * `NOTIFICATION_EMAIL_RECEIVER`

(Password needs to be generated - read [here](https://support.google.com/mail/?p=InvalidSecondFactor))

* `go build .` or `go run .`

The program updates itself and takes care of running the checks itself - you just need to run it. Once one of the mills will open the shop, it'll send you an email. Simple as that!
