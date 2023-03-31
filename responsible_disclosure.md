# Responsible Disclosure

When working with us according to this policy, you can expect us to:

* Work with you to understand and validate your report, including a timely initial response to the submission.
* Work to remediate discovered vulnerabilities in a timely manner.
* Recognize your contribution to improving our security if you are the first to report a unique vulnerability, and your report triggers a code or configuration change. Scope Any service on our domains *.indykite.com, *.indykite.id, *.indykite.me and *.indykite.one, but not services that are hosted by third parties. Examples of 3rd party services may include, but are not limited to, Akamai, Azure and GCloud. If you have any questions about scope, please contact us on [responsible-disclosure@indykite.com](mailto:responsible-disclosure@indykite.com)

When reporting vulnerabilities, please consider (1) attack scenario / exploitability, and (2) security impact of the bug.
‍
The following issues are considered out of scope:

* Clickjacking and UI redressing on pages with no sensitive actions.
* Cross-Site Request Forgery (CSRF) on unauthenticated forms or forms with no sensitive actions, e.g.: login/logout/search.
* Attacks requiring MITM or physical access to a user’s device.
* Previously known vulnerable libraries without a working Proof of Concept.
* Comma Separated Values (CSV) injection without demonstrating a vulnerability.
* Missing best practices in SSL/TLS configuration.
* Any activity that could lead to the disruption of our service (DoS).
* Content spoofing and text injection issues without showing an attack vector/without being able to modify HTML/CSS.
* Rate limiting or bruteforce issues on non-authentication endpoints.
* Missing best practices in Content Security Policy.
* Missing HttpOnly or Secure flags on cookies that are not sensitive (e.g. missing flags on authentication cookies are in scope).
* Missing email best practices (Invalid, incomplete or missing SPF/DKIM/DMARC records, etc.).
* Vulnerabilities only affecting users of outdated or unpatched browsers [Less than 2 stable versions behind the latest released stable version].
* Software version disclosure / Banner identification issues / Descriptive error messages or headers (e.g. stack traces, application or server errors).
* Public Zero-day vulnerabilities that have had an official patch for less than 1 month will be awarded on a case by case basis.
* Tabnabbing.
* Open redirect – unless an additional security impact can be demonstrated Issues that require unlikely user interaction.


# Rewards

We do not offer money or swag as rewards, but we will give you a place in our Security Hall of Fame. Reports need to meet the following criteria to qualify:

* Severity medium or above (severity is to be determined by IndyKite)
* Our security team is able to reproduce and verify the issue
* The issue is not previously known or already reported


# How to contact us

Our official communication channel is via email to:
[responsible-disclosure@indykite.com](mailto:responsible-disclosure@indykite.com)

# Safe Harbour
When conducting vulnerability research according to this policy, we consider this research conducted under this policy to be:

* Authorized in view of any applicable anti-hacking laws, and we will not initiate or support legal action against you for accidental, good-faith violations of this policy.
* Authorized in view of relevant anti-circumvention laws, and we will not bring a claim against you for circumvention of technology controls.
* Exempt from restrictions in our Acceptable Usage Policy that would interfere with conducting security research, and we waive those restrictions on a limited basis.
* Lawful, helpful to the overall security of the Internet, and conducted in good faith.
You are expected, as always, to comply with all applicable laws. If legal action is initiated by a third party against you and you have complied with this policy, we will take steps to make it known that your actions were conducted in compliance with this policy.

If at any time you have concerns or are uncertain whether your security research is consistent with this policy, please submit a report through one of our official channels before going any further.

# Hall of Fame
