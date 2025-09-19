package follow

import "news-feed/internal/middleware"

type Service interface {
	FollowUser(followerID, followeeID uint) error
	UnfollowUser(followerID, followeeID uint) error
	GetFollowingIDs(followerID uint) ([]uint, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FollowUser(followerID, followeeID uint) error {
	if followerID == followeeID {
		return middleware.BadRequest("cannot follow yourself")
	}

	isFollowing, err := s.repo.IsFollowing(followerID, followeeID)
	if err != nil {
		return err
	}
	if isFollowing {
		return middleware.Conflict("already following this user")
	}

	return s.repo.Follow(&Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
}

func (s *service) UnfollowUser(followerID, followeeID uint) error {
	return s.repo.Unfollow(followerID, followeeID)
}

func (s *service) GetFollowingIDs(followerID uint) ([]uint, error) {
	return s.repo.GetFollowingIDs(followerID)
}
