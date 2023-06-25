package cmd

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	addToPlaylistName string
)

var (
	trackID string
)

var (
	addTrackID                 string
	addTrackByIDToPlaylistName string
)

var (
	addTrackName                 string
	addTrackByNameToPlaylistName string
)

var (
	rmTrackName             string
	rmTrackFromPlaylistName string
)

var (
	newPlaylistName string
)

var (
	delPlaylistName string
)

var (
	listPlaylistTracksName string
)

func newCurrentTrackCmd() *cobra.Command {
	nowCmd := &cobra.Command{
		Use:   "now",
		Short: "Displays the currently playing track",
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayCurrentTrack(cmd, args)
		},
	}
	return nowCmd
}

func newShowTrackCmd() *cobra.Command {
	addToCmd := &cobra.Command{
		Use:   "show --tid [TRACK_ID]",
		Short: "Display information about a track by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayTrackById(cmd, args)
		},
	}
	addToCmd.Flags().StringVar(&trackID, "tid", "", "Id of track to display.")

	return addToCmd
}

func newAddToPlaylistCmd() *cobra.Command {
	addToCmd := &cobra.Command{
		Use:   "ato --p [PLAYLIST_NAME]",
		Short: "Add currently playing track to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTo(cmd, args)
		},
	}
	addToCmd.Flags().StringVar(&addToPlaylistName, "p", "", "Add current track to specified playlist.")

	return addToCmd
}

func newAddTrackByIDToPlaylistCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "aid --tid [TRACK_ID] --p [PLAYLIST_NAME]",
		Short: "Add track by ID to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTrackByIDToPlaylist(cmd, args)
		},
	}
	addCmd.Flags().StringVar(&addTrackID, "tid", "", "Id of track to add to playlist.")
	addCmd.Flags().StringVar(&addTrackByIDToPlaylistName, "p", "", "Name of playlist to add track to.")
	return addCmd
}

func newAddTrackByNameToPlaylistCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add --t [TRACK_NAME] --p [PLAYLIST_NAME]",
		Short: "Add track by name to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTrackByNameToPlaylist(cmd, args)
		},
	}
	addCmd.Flags().StringVar(&addTrackName, "t", "", "Name of track to add to playlist.")
	addCmd.Flags().StringVar(&addTrackByNameToPlaylistName, "p", "", "Name of playlist to add track to.")
	return addCmd
}

func newRemoveTrackFromPlaylistCmd() *cobra.Command {
	rmCmd := &cobra.Command{
		Use:   "rm --t [TRACK_NAME] --p [PLAYLIST_NAME]",
		Short: "Remove track from playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rmTrackByNameFromPlaylist(cmd, args)
		},
	}
	rmCmd.Flags().StringVar(&rmTrackName, "t", "", "Name of track to remove.")
	rmCmd.Flags().StringVar(&rmTrackFromPlaylistName, "p", "", "Name of playlist to remove track from.")
	return rmCmd
}

func newListPlaylistsCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "playlists",
		Short: "Show all playlists",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPlaylists(cmd, args)
		},
	}
	return newCmd
}

func newCreatePlaylistCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new --p [PLAYLIST_NAME]",
		Short: "Create new playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newPlaylist(cmd, args)
		},
	}
	newCmd.Flags().StringVar(&newPlaylistName, "p", "", "Name of new playlist.")
	return newCmd
}

func newDeletePlaylistCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "del --p [PLAYLIST_NAME]",
		Short: "Delete a playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deletePlaylist(cmd, args)
		},
	}
	deleteCmd.Flags().StringVar(&delPlaylistName, "p", "", "Name of playlist to delete.")
	return deleteCmd
}

func newListPlaylistTracksCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list --p [PLAYLIST_NAME]",
		Short: "List tracks in playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTracksFromPlaylist(cmd, args)
		},
	}
	listCmd.Flags().StringVar(&listPlaylistTracksName, "p", "", "Name of playlist to list tracks from.")
	return listCmd
}

func displayTrack(track *spotify.FullTrack) error {
	// format and display
	var data [][]interface{}
	item := []string{
		string(track.ID),
		track.Name,
		track.Album.Name,
		track.Artists[0].Name,
		(time.Duration(track.Duration) * time.Millisecond).Truncate(time.Second).String(),
		strconv.Itoa(track.Popularity),
		strconv.FormatBool(track.Explicit),
		track.PreviewURL,
	}
	row := make([]interface{}, len(item))
	for i, d := range item {
		row[i] = d
	}
	data = append(data, row)
	printSimple([]string{"ID", "Name", "Album", "Artist", "Duration", "Popularity", "Explicit", "Preview"}, data)
	return nil
}

func displayTrackById(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the track (check for existence)
	track, err := client.GetTrack(spotify.ID(trackID))
	if err != nil {
		return err
	}

	err = displayTrack(track)
	if err != nil {
		return err
	}

	return nil
}

func displayCurrentTrack(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get current playing song
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}

	err = displayTrack(playing.Item)
	if err != nil {
		return err
	}

	return nil
}

func addTo(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get my playlists
	pl, err := getPlaylistByName(addToPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// get current playing song
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}
	fmt.Println("Track: ", playing.Item.Name)

	// add track to playlist
	_, err = client.AddTracksToPlaylist(pl.ID, playing.Item.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Added track \"%s\" to playlist \"%s\".\n", playing.Item.Name, pl.Name)
	return nil
}

func listPlaylists(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get all playlists for the user
	playlists, err := getPlaylists()
	if err != nil {
		return err
	}

	// format resulting data
	var data [][]interface{}
	if playlists.Playlists != nil {
		for _, item := range playlists.Playlists {
			track := []string{
				string(item.ID),
				item.Name,
				item.Owner.DisplayName,
				strconv.FormatBool(item.IsPublic),
				strconv.FormatBool(item.Collaborative),
				strconv.FormatUint(uint64(item.Tracks.Total), 10)}
			row := make([]interface{}, len(track))
			for i, d := range track {
				row[i] = d
			}
			data = append(data, row)
		}
	}

	// pretty print track results
	printSimple([]string{"ID", "Name", "Owner", "Public", "Collaborative", "Tracks"}, data)
	return nil
}

func newPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// create new playlist
	playlist, err := client.CreatePlaylistForUser(user.ID, newPlaylistName, "", true)
	if err != nil {
		return err
	}
	fmt.Println("Created public playlist: ", playlist.Name)
	return nil
}

func deletePlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist
	pl, err := getPlaylistByName(delPlaylistName)
	if err != nil {
		return err
	}

	// unfollow and return
	// TODO: delete != unfollow?
	return client.UnfollowPlaylist(spotify.ID(user.ID), pl.ID)
}

func addTrackByIDToPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(addTrackByIDToPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// get the track (check for existence)
	tr, err := client.GetTrack(spotify.ID(addTrackID))
	if err != nil {
		return err
	}
	fmt.Println("Track: ", tr.Name)

	// add track to playlist
	_, err = client.AddTracksToPlaylist(pl.ID, tr.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Added track \"%s\" to playlist \"%s\".\n", tr.Name, pl.Name)
	return nil
}

func addTrackByNameToPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(addTrackByNameToPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// Search for the track
	results, err := client.Search(addTrackName, spotify.SearchTypeTrack)
	if err != nil {
		return err
	}

	// add most popular track to playlist from results
	if results.Tracks != nil {
		tracks := results.Tracks.Tracks[:]
		sort.Slice(tracks, func(i, j int) bool { return tracks[i].Popularity > tracks[j].Popularity })
		fmt.Println("Track: ", tracks[0].Name)

		// add track to playlist
		_, err = client.AddTracksToPlaylist(pl.ID, tracks[0].ID)
		if err != nil {
			return err
		}
		fmt.Printf("Added track \"%s\" to playlist \"%s\".\n", tracks[0].Name, pl.Name)
	} else {
		fmt.Printf("Track %s not found.\n", addTrackName)
	}
	return nil
}

func rmTrackByNameFromPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(rmTrackFromPlaylistName)
	if err != nil {
		return err
	}

	// get track in playlist and validate existence
	var matchedTrack spotify.SimpleTrack
	ptracks, err := client.GetPlaylistTracks(pl.ID)
	if err != nil {
		return err
	}
	for _, t := range ptracks.Tracks {
		if rmTrackName == t.Track.SimpleTrack.Name {
			matchedTrack = t.Track.SimpleTrack
			break
		}
	}
	if reflect.DeepEqual(matchedTrack, spotify.SimpleTrack{}) {
		return fmt.Errorf("track %s not found in playlist %s", rmTrackName, rmTrackFromPlaylistName)
	}
	fmt.Println("Track: ", matchedTrack.Name)

	// remove track from playlist
	_, err = client.RemoveTracksFromPlaylist(pl.ID, matchedTrack.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Removed track \"%s\" from playlist \"%s\".\n", matchedTrack.Name, rmTrackFromPlaylistName)
	return nil
}

func listTracksFromPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	pl, err := getPlaylistByName(listPlaylistTracksName)
	if err != nil {
		return err
	}

	// get tracks from playlist
	tracks, err := client.GetPlaylistTracks(pl.ID)
	if err != nil {
		return err
	}

	// format resulting data
	var data [][]interface{}
	if tracks.Tracks != nil {
		for _, item := range tracks.Tracks {
			track := []string{
				string(item.Track.ID),
				item.Track.Name,
				item.Track.Album.Name,
				item.Track.Artists[0].Name,
				strconv.Itoa(item.Track.Popularity)}
			row := make([]interface{}, len(track))
			for i, d := range track {
				row[i] = d
			}
			data = append(data, row)
		}
	}

	// pretty print track results
	printSimple([]string{"ID", "Name", "Album", "Artist", "Popularity"}, data)
	return nil
}

func getPlaylists() (*spotify.SimplePlaylistPage, error) {
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return &(spotify.SimplePlaylistPage{}), err
	}

	return playlists, nil
}

func getPlaylistByName(playlistName string) (spotify.SimplePlaylist, error) {
	// get current user's playlists
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return spotify.SimplePlaylist{}, err
	}

	// match by name
	var matchPlaylist spotify.SimplePlaylist
	for _, p := range playlists.Playlists {
		if playlistName == p.Name {
			matchPlaylist = p
			break
		}
	}

	// check if found and return
	if reflect.DeepEqual(matchPlaylist, spotify.SimplePlaylist{}) {
		return spotify.SimplePlaylist{}, fmt.Errorf("playlist not found: %s", playlistName)
	}
	return matchPlaylist, nil
}
