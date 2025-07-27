# Security Policy

## 🔒 Supported Versions

We actively support the following versions of MCP RAPI with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | ✅ Yes             |
| < 1.0   | ❌ No              |

## 🚨 Reporting a Vulnerability

We take the security of MCP RAPI seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### 📧 How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please send an email to: **security@mcprapi.dev** (or create a private security advisory on GitHub)

Include the following information in your report:
- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

### 🕐 Response Timeline

- **Initial Response**: Within 48 hours of receiving your report
- **Status Update**: Within 7 days with a more detailed response
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days

### 🎯 What to Expect

After submitting a report, you can expect:

1. **Acknowledgment**: We'll confirm receipt of your vulnerability report
2. **Investigation**: We'll investigate and validate the reported vulnerability
3. **Resolution**: We'll work on a fix and coordinate the release
4. **Disclosure**: We'll publicly disclose the vulnerability after a fix is available

### 🏆 Recognition

We believe in recognizing security researchers who help keep our users safe. With your permission, we'll:
- Credit you in our security advisories
- Add you to our security researchers hall of fame
- Provide a reference letter if requested

## 🛡️ Security Best Practices

### For Users

1. **Keep Updated**: Always use the latest version of MCP RAPI
2. **Secure Configuration**: Follow our security configuration guidelines
3. **Environment Variables**: Never commit sensitive environment variables
4. **HTTPS**: Always use HTTPS in production
5. **Database Security**: Secure your database connections and credentials
6. **Regular Audits**: Regularly audit your permissions and access controls

### For Developers

1. **Input Validation**: Always validate and sanitize user inputs
2. **Authentication**: Implement proper authentication and authorization
3. **Secrets Management**: Never hardcode secrets or credentials
4. **Dependencies**: Regularly update dependencies and scan for vulnerabilities
5. **Code Review**: All code changes must go through security review
6. **Testing**: Include security testing in your development process

## 🔍 Security Features

MCP RAPI includes several built-in security features:

- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control (RBAC)**: Fine-grained permission management
- **Rate Limiting**: Protection against abuse and DoS attacks
- **Input Validation**: Comprehensive input sanitization
- **CORS Protection**: Configurable cross-origin resource sharing
- **SQL Injection Prevention**: Parameterized queries and ORM protection
- **XSS Protection**: Output encoding and CSP headers
- **Secure Headers**: Security-focused HTTP headers

## 📚 Security Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)
- [Vue.js Security Guide](https://vuejs.org/guide/best-practices/security.html)
- [Docker Security Best Practices](https://docs.docker.com/engine/security/)

## 🔄 Security Updates

Security updates will be:
- Released as patch versions (e.g., 1.0.1, 1.0.2)
- Documented in our changelog with security advisory references
- Announced through our security mailing list
- Tagged with appropriate severity levels (Critical, High, Medium, Low)

## 📞 Contact

For any security-related questions or concerns:
- Email: security@mcprapi.dev
- Security Advisory: [GitHub Security Advisories](https://github.com/chuangyeshuo/mcprapi/security/advisories)

---

**Thank you for helping keep MCP RAPI and our users safe!** 🙏