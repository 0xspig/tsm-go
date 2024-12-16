---
title: "Made a couple scripts"
date: 2023-05-08T01:22:32-04:00
category: blog
draft: false
sketch: false
tags: 
- technology
- updates

---

## Build and Deploy Script

    #!/bin/bash
    
    (cd /home/tyler/www/tsm ;
    hugo ;
    echo "uploading";
    rsync -rtvzP ./public/* root@tsmckee.com:/var/www/tsm;
    )

## New Blog Script

    #!/bin/bash
    
    NAME=$1
    
    (cd /home/tyler/www/tsm/;
    hugo new blog/$NAME;
    nvim content/blog/$NAME;
    )


It's kind of amazing how easy it is to program on linux. Or rather, its amazing how difficult windows makes it. 
