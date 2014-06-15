/*
 * print.c - common print support functions for lsof
 */


/*
 * Copyright 1994 Purdue Research Foundation, West Lafayette, Indiana
 * 47907.  All rights reserved.
 *
 * Written by Victor A. Abell
 *
 * This software is not subject to any license of the American Telephone
 * and Telegraph Company or the Regents of the University of California.
 *
 * Permission is granted to anyone to use this software for any purpose on
 * any computer system, and to alter it and redistribute it freely, subject
 * to the following restrictions:
 *
 * 1. Neither the authors nor Purdue University are responsible for any
 *    consequences of the use of this software.
 *
 * 2. The origin of this software must not be misrepresented, either by
 *    explicit claim or by omission.  Credit to the authors and Purdue
 *    University must appear in documentation and sources.
 *
 * 3. Altered versions must be plainly marked as such, and must not be
 *    misrepresented as being the original software.
 *
 * 4. This notice may not be removed or altered.
 */

#ifndef lint
static char copyright[] =
"@(#) Copyright 1994 Purdue Research Foundation.\nAll rights reserved.\n";
static char *rcsid = "$Id: print.c,v 1.55 2013/01/02 17:14:59 abe Exp $";
#endif


#include "lsof.h"


/*
 * Local definitions, structures and function prototypes
 */



struct porttab {
	int port;
	MALLOC_S nl;			/* name length (excluding '\0') */
	int ss;				/* service name status, 0 = lookup not
					 * yet performed */
	char *name;
	struct porttab *next;
};


#if	defined(HASNORPC_H)
static struct porttab **Pth[2] = { NULL, NULL };
						/* port hash buckets:
						 * Pth[0] for TCP service names
						 * Pth[1] for UDP service names
						 */
#else	/* !defined(HASNORPC_H) */
static struct porttab **Pth[4] = { NULL, NULL, NULL, NULL };
						/* port hash buckets:
						 * Pth[0] for TCP service names
						 * Pth[1] for UDP service names
						 * Pth[2] for TCP portmap info
						 * Pth[3] for UDP portmap info
						 */
#endif	/* defined(HASNORPC_H) */


#if	!defined(HASNORPC_H)
_PROTOTYPE(static void fill_portmap,(void));
_PROTOTYPE(static void update_portmap,(struct porttab *pt, char *pn));
#endif	/* !defined(HASNORPC_H) */

_PROTOTYPE(static void fill_porttab,(void));
_PROTOTYPE(static char *lkup_svcnam,(int h, int p, int pr, int ss));
_PROTOTYPE(static int printinaddr,(void));


/*
 * endnm() - locate end of Namech
 */

char *
endnm(sz)
	size_t *sz;			/* returned remaining size */
{
	register char *s;
	register size_t tsz;

	for (s = Namech, tsz = Namechl; *s; s++, tsz--)
		;
	*sz = tsz;
	return(s);
}


#if !defined(HASNORPC_H)
/*
 * fill_portmap() -- fill the RPC portmap program name table via a conversation
 *		     with the portmapper
 *
 * The following copyright notice acknowledges that this function was adapted
 * from getrpcportnam() of the source code of the OpenBSD netstat program.
 */

/*
* Copyright (c) 1983, 1988, 1993
*      The Regents of the University of California.  All rights reserved.
*
* Redistribution and use in source and binary forms, with or without
* modification, are permitted provided that the following conditions
* are met:
* 1. Redistributions of source code must retain the above copyright
*    notice, this list of conditions and the following disclaimer.
* 2. Redistributions in binary form must reproduce the above copyright
*    notice, this list of conditions and the following disclaimer in the
*    documentation and/or other materials provided with the distribution.
* 3. All advertising materials mentioning features or use of this software
*    must display the following acknowledgement:
*      This product includes software developed by the University of
*      California, Berkeley and its contributors.
* 4. Neither the name of the University nor the names of its contributors
*    may be used to endorse or promote products derived from this software
*    without specific prior written permission.
*
* THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
* ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
* IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
* ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
* FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
* DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
* OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
* HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
* LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
* OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
* SUCH DAMAGE.
*/

static void
fill_portmap()
{
    // stripped
}
#endif	/* !defined(HASNORPC_H) */


/*
 * fill_porttab() -- fill the TCP and UDP service name port table with a
 *		     getservent() scan
 */

static void
fill_porttab()
{
    // stripped
}


/*
 * gethostnm() - get host name
 */

char *
gethostnm(ia, af)
	unsigned char *ia;		/* Internet address */
	int af;				/* address family -- e.g., AF_INET
					 * or AF_INET6 */
{
    // stripped
}


/*
 * lkup_svcnam() - look up service name for port
 */

static char *
lkup_svcnam(h, p, pr, ss)
	int h;				/* porttab hash index */
	int p;				/* port number */
	int pr;				/* protocol: 0 = TCP, 1 = UDP */
	int ss;				/* search status: 1 = Pth[pr][h]
					 *		  already searched */
{
    // stripped
}


/*
 * print_file() - print file
 */

void
print_file()
{
}


/*
 * printinaddr() - print Internet addresses
 */

static int
printinaddr()
{
    // stripped
}


/*
 * print_init() - initialize for printing
 */

void
print_init()
{
}


#if	!defined(HASPRIVPRIPP)
/*
 * printiproto() - print Internet protocol name
 */

void
printiproto(p)
	int p;				/* protocol number */
{
}
#endif	/* !defined(HASPRIVPRIPP) */


/*
 * printname() - print output name field
 */

void
printname(nl)
	int nl;				/* NL status */
{
    // stripped
}


/*
 * printrawaddr() - print raw socket address
 */

void
printrawaddr(sa)
	struct sockaddr *sa;		/* socket address */
{
    // stripped
}


/*
 * printsockty() - print socket type
 */

char *
printsockty(ty)
	int ty;				/* socket type -- e.g., from so_type */
{
    // stripped
}


/*
 * printuid() - print User ID or login name
 */

char *
printuid(uid, ty)
	UID_ARG uid;			/* User IDentification number */
	int *ty;			/* returned UID type pointer (NULL
					 * (if none wanted).  If non-NULL
					 * then: *ty = 0 = login name
					 *	     = 1 = UID number */
{
// stripped
}


/*
 * printunkaf() - print unknown address family
 */

void
printunkaf(fam, ty)
	int fam;			/* unknown address family */
	int ty;				/* output type: 0 = terse; 1 = full */
{
    // stripped
}
