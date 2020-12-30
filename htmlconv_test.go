package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testWpoData = `
<div id="readability-page-1" class="page"><div><div><article dir="ltr"><div><p>WASHINGTON (Reuters) - Small drones will be allowed to fly over people and at night in the United States, the Federal Aviation Administration (FAA) said on Monday, a significant step toward their use for widespread commercial deliveries.</p><p>The FAA said its long-awaited rules for the drones, also known as unmanned aerial vehicles, will address security concerns by requiring remote identification technology in most cases to enable their identification from the ground.</p><p>Previously, small drone operations over people were limited to operations over people who were directly participating in the operation, located under a covered structure, or inside a stationary vehicle - unless operators had obtained a waiver from the FAA.</p><p>The rules will take effect 60 days after publication in the federal register in January. Drone manufacturers will have 18 months to begin producing drones with Remote ID, and operators will have an additional year to provide Remote ID.</p><p>There are other, more complicated rules that allow for operations at night and over people for larger drones in some cases.</p><p>“The new rules make way for the further integration of drones into our airspace by addressing safety and security concerns,” FAA Administrator Steve Dickson said. “They get us closer to the day when we will more routinely see drone operations such as the delivery of packages.”</p><p>Companies have been racing to create drone fleets to speed deliveries. The United States has over 1.7 million drone registrations and 203,000 FAA-certificated remote pilots.</p><p>For at-night operations, the FAA said drones must be equipped with anti-collision lights. The final rules allow operations over moving vehicles in some circumstances.</p><p>Remote ID is required for all drones weighing 0.55 lb (0.25 kg) or more, but is required for smaller drones under certain circumstances like flights over open-air assemblies.</p><p>The new rules eliminate requirements that drones be connected to the internet to transmit location data but do that they broadcast remote ID messages via radio frequency broadcast. Without the change, drone use could have been barred from use in areas without internet access.</p><p>The Association for Unmanned Vehicle Systems International said Remote ID will function as “a digital license plate for drones ... that will enable more complex operations” while operations at night and over people “are important steps towards enabling integration of drones into our national airspace.”</p><p>One change, since the rules were first proposed in 2019, requires that small drones not have any exposed rotating parts that would lacerate human skin.</p><p>United Parcel Service Inc said in October 2019 that it won the government’s first full approval to operate a drone airline.</p><p>Last year, Alphabet’s Wing, a sister unit of search engine Google, was the first company to get U.S. air carrier certification for a single-pilot drone operation.</p><p>In August, Amazon.com Inc’s drone service received federal approval allowing the retailer to begin testing commercial deliveries through its drone fleet.</p><p>Walmart Inc said in September it would run a pilot project for delivery of grocery and household products through automated drones but acknowledged “it will be some time before we see millions of packages delivered via drone.”</p><div><p>Reporting by David Shepardson; Editing by Nick Zieminski and Howard Goller</p></div></div></article></div></div><p><span data-name="for-phone-only">for-phone-only</span><span data-name="for-tablet-portrait-up">for-tablet-portrait-up</span><span data-name="for-tablet-landscape-up">for-tablet-landscape-up</span><span data-name="for-desktop-up">for-desktop-up</span><span data-name="for-wide-desktop-up">for-wide-desktop-up</span></p></div>
`

const testSimpleData = `<p>it <b>works</b>!</p><p>
<a href="https://ya.ru\dumb">link <b>text</b></a>
</p>
`

func TestWpo(t *testing.T) {
	so := assert.New(t)

	got, err := htmlToTex(testSimpleData)
	fmt.Println(got)
	so.Nil(err)
}
