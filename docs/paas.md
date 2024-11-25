Rails recently released [Rails 8.0: No PaaS Required](https://rubyonrails.org/2024/11/7/rails-8-no-paas-required).  Fly.io seeks to redefine PaaS.  Let's start by defining a few terms.

# IaaS: Infrastructure as a Service

Essentially, this is Virtual Private Servers (VPS) that you can rent.  Some of the bigger players include: Amazon Elastic Compute Cloud (EC2) and Google Compute Engine.  Rails
encourages Hetzner and Digital Ocean Droplets.

Essentially, what you can get for a few Euros a month is root access to a virtual machine with an operating system installed and an IP address.  You request one and it is ready in minutes.

[Fly Machines are fast-launching VMs](https://fly.io/docs/machines/).  They cost a few pennies per hour, and are ready in milliseconds.

There are 730 hours per month, so if you do the math, a Fly Machine can cost more than a Hetzner VPS, but a Fly machine does more, and are configured by default auto stop when not in use, so it typically costs less.

# PaaS: Platform as a Service

Installing, configuring, updating, and backing up virtual machines is a bit of a chore.

Companies like Heroku, Render, and Railway lighten this load by providing pre-configured
platforms for your application, including databases and other services.  In DHH's
[Rails World 2024 Opening Keynote](https://www.youtube.com/watch?v=-cEn_83zRFw), he pointed out that many of these are build on third party Clouds and bill you for both
set of services.

Rails 8 provides a Dockerfile and a tool named Kamal that will configure and deploy your application to a standard VPS.  You are still responsible for firewalls, load balancers, backups, databases, and more.

Fly launch can deploy your application using the exact same Dockerfiles that Kamal uses.
We provide the firewalls, load balancers, backups, databases, and more.

# Sidebar: Rails 8 doesn't require a PaaS.  

Truth is, Rails never has required a PaaS.  Don't believe me?  Here's a blog post from 2008: [Myth #1: Rails is hard to deploy](https://dhh.dk/posts/30-myth-1-rails-is-hard-to-deploy), which ends with:

> In conclusion, Rails is no longer hard to deploy. Phusion Passenger has made it ridiculously easy.

I am confident that there are still many people out there who are happily deploying their
Rails 6 applications using Capistrano and Phusion Passenger.  Replacing these components
with Kamal and Thruster doesn't materially change the value proposition.

There are those that benefit from provisioning their own machines.  There are those who
would rather offload this responsibility to others.

The choice is yours to make.

# SaaS: Software as a Service

This is where you come in.  Fly.io is a Developer-Focused Public Cloud.  You write
software.  You write software, and want to make it available online as a service.  That's
our specialty.

Wikipedia on [Software as a service](https://en.wikipedia.org/wiki/Software_as_a_service):
  * by 2023 was the main form of software application deployment
  * SaaS architectures are typically multi-tenant; usually they share resources between clients for efficiency, but sometimes they offer a siloed environment for an additional fee.

fly-atc is a SaaS toolkit for converting a personal application into a efficient, siloed, multi-tenant application, where each user of your application can be assigned a dedicated virtual machine.
