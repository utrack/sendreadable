package tpl

import (
	"io"
	"text/template"
)

var tpl = template.Must(template.New("").Delims("$$", "$$").Parse(text))

type Request struct {
	Title         string
	Author        string
	URL           string
	SourceName    string
	AvgTimeString string

	Content string

	FontPath string
}

func Render(r Request, w io.Writer) error {
	return tpl.Execute(w, r)
}

const text = `
% universal settings
\documentclass[a4paper,12pt,oneside]{article}
\usepackage{anyfontsize}
%\usepackage[utf8x]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{fontspec}
\usepackage[main=english,russian]{babel}
\setsansfont[Path = $$ .FontPath $$/] {BasisGrotesquePro-Regular}
\setmainfont[Path = $$ .FontPath $$/, BoldFont={BasisGrotesquePro-Bold}, ItalicFont={BasisGrotesquePro-Italic}, BoldItalicFont={BasisGrotesquePro-BoldItalic}]{BasisGrotesquePro-Regular}

\usepackage{graphicx}
% \setmonofont[ Path = fonts/,  ] { }
\usepackage{lettrine}
\usepackage{enumitem}
\usepackage{hyperref}
\usepackage{titlesec}
\usepackage{xcolor}
% pdflatex only
% \usepackage[tracking=true, letterspace=50]{microtype}
\usepackage[left=0.75in, right=0.75in, top=1in, bottom=1in]{geometry}
\usepackage{setspace}
\usepackage{calc}
\usepackage{qrcode}
\usepackage{tabularx}
\usepackage{multirow}
\usepackage[export]{adjustbox} % loads also graphicx
\usepackage[font=small]{caption}
\usepackage{booktabs} % table hlines
\usepackage{ltablex}
\usepackage{eso-pic} % topright qrcode

\usepackage{etoolbox}
\AtBeginEnvironment{quote}{\singlespacing\small}

\title{$$ .Title $$}
\author{$$ .Author $$$$ if .SourceName $$ | \href{$$ .URL $$}{$$ .SourceName $$}$$ end $$}
\newcommand{\readingTime}{$$ .AvgTimeString $$}


% basic global settings
\linespread{1.4}
% pdflatex only
% \SetTracking{encoding=*,shape=sc}{50}
\setlist[itemize]{itemsep=0em}
\setlength{\parskip}{0.75em}

% qrcode at top right
\newcommand\AtPageUpperRight[1]{\AtPageUpperLeft{%
   \makebox[\paperwidth][r]{#1}}}

% title format
\makeatletter
\def\maketitle{\noindent{
\begin{flushleft}
 {\fontsize{26}{0}\selectfont\sffamily\bfseries\@title}\\
 \vspace{1em}
\end{flushleft}
$$ if ne .Author "" $${\@author}\\$$ end $$
{\small{\readingTime}}\\
\vspace{1em}
{\hrulefill}
  }
}
\makeatother

%\titleformat{\section}
%  {\normalfont\Large\bfseries}{\thesection}{1em}{}

\begin{document}

\AddToShipoutPictureBG*{%
  \AtPageUpperRight{\raisebox{-\height}{\frame{{\qrcode[hyperlink,level=Q,tight,height=3cm]{$$ .URL $$}}}}}}


\maketitle

$$ .Content $$

\end{document}
`
