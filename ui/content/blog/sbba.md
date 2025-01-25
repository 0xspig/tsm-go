---
title: Saw Blade Bevel Angles
date: 2025-01-23T17:13:20-04:00
category: blog
draft: true
image: 
imgalt: 
tags:
  - wood  
  - geometry
---
# Arris as the Ground Line
I've been trying to understand this technique for laying out a compound angle birdsmouth called [Four Hits of the Square](https://www.youtube.com/watch?v=stzrhTLiWFk).
Pat Moore demonstrates how to find the birdsmouth angle and saw blade bevel angles from just the plum and alignment angles on the timber.
This is faster than the way I did it in {my hip roof models}(hips.md).
The flipside of it's brevity that its less obvious what is happening than on paper.
The tl;dr on these techniques is that the arris of the timber is acting as the ground line in our full drawings.
If that doesn't make sense yet give me a couple minutes and I'll explain.

## SBBA
We'll start with the simpler trick for finding the saw blade bevel angle.
On a simple bevel cut, you can set the saw blade to your bevel angle and run it.
If you you cut at an angle, however, your bevel gets skewed.
Think of a flashlight on a wall.
If you hold the flashlight perpendicular to the wall you get a circle.
If you hold it to the side you get an elipse.
That's what happens to your bevel angle.

We'll start on paper with a Ground Line and three points representing the angle between the top of our timber and the true angle of the bevel cut or plumb cut.
Above GL is plan view (looking down on the timber) and below is elevation (looking at the side of the timber).

<img src="/static/images/blog/sbba/00_GL.jpg" width="80%">

The next thing we need is a line to follow with our saw - the alignment cut. 
Keep in mind here that these two lines form an acute angle.
The technique will not work if they dont and you'll need to roll your timber 90 degrees.
We're going to scribe a plane that is perpendicular to alignment cut as well as the horizontal plane.
HT is the intersection of this plane and the horizontal plane.
VT is the intersection of this plane and the vertical plane.

<img src="/static/images/blog/sbba/01_P.jpg" width="80%">

Now the fun part, we project the points of our original drawing onto our new plane.
The H components of our points are projected onto the horizontal trace of our plane (HT) and V components to the vertical trace (VT) - plan to plan, elevation to elevation.

<img src="/static/images/blog/sbba/02_P_Proj.jpg" width="80%">

The final step is rotating this plane into a view plane.
You can rotate the plane around either HT or VT but for this demonstration we'll use the vertical axis.
While a clockwise rotation is shorter, the result overlaps our existing drawing.
For the sake of clarity we'll rotate it counterclockwise, giving us a mirror image of our angle.

<img src="/static/images/blog/sbba/03_Fin.jpg" width="80%">

As you can see, the points have been scaled down along their horizontal axis resulting in a shallower angle.

## Faster
That demo got us the right answer but we can streamline it quite a bit.
For starters, we'll move the top edge of our model up to the Ground Line.
While we're at it we can place the H components on the Ground Line as well.

Next, consider that any point which intersects our projection plane will not be moved when the plane rotates.
This time when we scribe our projection plane we will do so through point C.

Now we know where C' will land after the rotation. 
Point A becomes redundant because that edge of the timber is just a ray from B along GL.
We don't have to worry about B's V component either because it is locked in GL.
The only point we have to project now is B's H component.
We scribe the alignment perpendicular, project Bh onto it, and rotate our plane into the view plane giving us B'.

The angle between GL, B' and Cv is our saw blade bevel angle.
To find an SBBA on site, just substitute the Ground Line of this process for an edge on your timber.

I'll make another post at some point when I figure out the birdsmouth cut.