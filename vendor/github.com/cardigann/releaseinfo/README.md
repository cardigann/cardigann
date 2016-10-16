
Release Info Parser
===================

A library for parsing media file names and paths and returning structured information about them. 

Ported from [Sonarr](https://github.com/Sonarr/Sonarr/commit/f2ccf948356404beceb1e5dea73f7db79cb64dd3)'s parser code. 

## Example

```go

result, err := releaseinfo.Parse("WEEDS.S03E01-06.DUAL.Bluray.AC3.-HELLYWOOD.avi")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%#v", result)

=> &releaseinfo.EpisodeInfo{
    SeriesTitle:"WEEDS",
    SeriesTitleInfo:releaseinfo.SeriesTitleInfo{
        Title:"WEEDS", 
        TitleWithoutYear:"WEEDS", 
        Year:0
    }, 
    Quality:releaseinfo.QualityModel{
        Quality:releaseinfo.Quality{Id:6, Name:"Bluray-720p"}, 
        Revision:0, 
        QualitySource:"name"
    }, 
    SeasonNumber:3, 
    EpisodeNumbers:[]int{1, 2, 3, 4, 5, 6}, 
    AbsoluteEpisodeNumbers:[]int{}, 
    AirDate:"", 
    Language:language.Tag{lang:0x9a, region:0x0, script:0x0, pVariant:0x0, pExt:0x0, str:""}, 
    FullSeason:false, 
    Special:false, 
    ReleaseGroup:"HELLYWOOD", 
    ReleaseHash:""
}
```
