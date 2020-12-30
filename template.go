package main

const tplLatex = `
% universal settings
\documentclass[a4paper,12pt,oneside]{article}
\usepackage{anyfontsize}
%\usepackage[utf8x]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{fontspec}
\setmainfont[ Path = /tmp/sendreadable/fonts/] {BasisGrotesquePro-Regular}
\setsansfont[ Path = /tmp/sendreadable/fonts/] {BasisGrotesquePro-Regular}

\usepackage{graphicx}
% \setmonofont[ Path = fonts/,  ] { }
\usepackage{lettrine}
\usepackage{enumitem}
\usepackage{hyperref}
\usepackage{titlesec}
\usepackage{xcolor}
% pdflatex only
% \usepackage[tracking=true, letterspace=50]{microtype}
\usepackage[left=1.25in, right=1.25in, top=1in, bottom=1.25in]{geometry}
\usepackage{setspace}
\usepackage{calc}
\usepackage{qrcode}
\usepackage{wrapfig}
\usepackage{supertabular}
\usepackage{multirow}


\title{$$ .Title $$}
\author{$$ .Author $$ | \href{$$ .URL $$}{$$ .Source $$}}
\newcommand{\readingTime}{$$ .AvgTimeString $$}


% basic global settings
\linespread{1.5}
% pdflatex only
% \SetTracking{encoding=*,shape=sc}{50}
\setlist[itemize]{itemsep=0em}
\setlength{\parskip}{0.75em}

% title format
\makeatletter
\def\maketitle{\noindent{
    \begin{flushleft}
        {\fontsize{26}{0}\selectfont\sffamily\bfseries\@title}\\\vspace{1em}%
        \begin{tabular}{l r}

          $$ if ne .Author "" $${\@author}$$ end $$\vspace{0.3em} & \multirow{2}{*}{\hspace{2cm}\qrcode[hyperlink,level=Q,tight]{$$ .URL $$}}\\%
          {\small{\readingTime}} &
        \end{tabular}\\
        \vspace{2em}
        {\hrulefill}%
    \end{flushleft}
  }
}
\makeatother

\titleformat{\section}
  {\normalfont\sffamily\Large\AlegreyaSansExtraBold}
  {\thesection}{}{}

\begin{document}

\maketitle

$$ if ne .Image "" $$
\begin{center}
  \makebox[\textwidth]{\includegraphics[width=\textwidth]{$$ .Image $$}}
\end{center}
$$ end $$

$$ .Content $$

\end{document}
`
