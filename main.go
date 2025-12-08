package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=ca05be6623084b93b8b01c9e1200e0c3; _kpiid=8bc3c90e0a818233479d4791be94e941; _kau=43e5aed7b4e3bc7363326df25da0df67b8762248644dc28bf7bf9892a1d8e4b4b17e8983b6262e529ca249525dad6e42c1f586a378138a86fe13ebd51f07fe96f0df78152ab60dbbff3627864a644f4902b5b6af826ecbe907216a3c0a8d596996c7458be8cebd9b96e9854d514e2e7ea3f108153ce281a81e3e65d73636353638313638333830313531393734323539373136343032323132343830ec9f75f14d8f7d231bb23fecdd988f61; _kawlt=Yj5KpONMLdSd9IfmuSO8aH2dPgwIbCXZ4X8BEUewMSAKooTSgvNaPkK3dxM57DVptk6GCil3lT7iTeZ23ep8_EjFWhbiGTWcrxXeOA0Cc5WF_dFXKKWruJSVj4GB83wo; _kawltea=1765289441; _karmt=RzZI5Q_n7_KeI4px4y94D5hIXCu8wJwhSxDQk1E4Ou23bIRxZDxj1WudInVcXqqf; _karmtea=1765300241; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kpawbat_e=J8tXjo4DRmJvgNjkZttQ68zHHk8s43fOrxN%2FDojlW01rqQh8rqYluCr4n%2FwGz9Jrt4uQY2g9%2BgG4AIFVIME1J0Meo8SyzFx6A86VSSXhEQU%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..V7gJs5xnlRjvN_Xt.cUzlw58UF67grUbvNEqJl0TjB_P5ME5NsWR0lJZOOtmgNY8ulJrOqm252VYAJRd6xTpZ5QgduTz7Ourb-ObUZSUPYHyFHO0ZWW5YsgKJnFU6SsWOlCBI67Oiif6RJNdBpwxX3ksS0IkSa3jmIksgsqGP3NdkGRIXoINFPhsquKOa6TBo-nY5IIZNUe9mtnDQbkHSRfufxjqEVeUCVxl4hZdD1bZiRf0BpqFjxH5N1OOKGEh2varbjyQ5dqp3aoP7T3iKHgY_zTiD4I-nh9A8tsSd4pWwg9NfkQkEENwDCKSwAQ3ngV16XovKnGrEhYaPdreLNrBsppSVEU30TUuX1dHjYOonAq2VE4UmxaLP5T-3Aokpha8nbMigwzon18Z2bZF_RyKV8PIV3z_O1aHtX4BYOu3EYAOOqAikxt3jDwq4jaqTNdn5Ke9qubmTfgx6gAjn1XiKptEOoNpWpeW7XdN345XWml_TngywU3C32ct7D_WJCo_HjGxK5jJLxPEM7rpPvBplZiYJzykOOk-BylEYTU8szPC2kjKePSWEstwOmQJsEL4cXLTUK8JGRCoZqr1uhujlRVzXma6-0QtwbZADXm13S9FutGZ6sFJ5TvVNe9UWFt1ipoK5ZMTO_HEKUimQc-y8WwE8slhBrGSIq3u3WaNsnKX6wJuEzMVJrPToaPIhp7ZJ6XYiGGN-V1XfB0s7o-3mHYnd4u5ezBLc5v-vvow5f5UVnleCOlTATfK-KyjHrjy1uCnOOuZeOV0-LatDEmgmdq3Wcx9J70g_gcTvB6M6HMTdMmLIvVUGQB7gAyhqZ8K1Hrb_fqV6fnfH7YMrWls12RB5ZwRTYuwKlaNEAAfgtr6Mk9mFHA6JaOpkZIciVq-qouZX_OMGsY2yExlxeLJcItkJPn_w1xJXJtqtsk0MqBTCjtp4lohemY_6TpGYtAWqXRtGqp8sosQri8AQ0QwDO8DkHoZjwZNBUvxFXFyawoKHYRjr9nYEgXLeo31e1B0dnRIwtO6ZFqs9vll_3vVd38ga9cdlZmJE52gBqGvjR_5xC2ZZwOJDi-gZTbJ6ATynP46dnUnkdRi0Y5xLyT8n4QdWQR6F3190bb4e0hKjl8Fhvbeyu9am_QF8m5GlvyPki6anIHXyhNUSQeWBsvIZsWxmabAM00JjcreRZ6Dx638DrbGK-GC-WpLkXfPr8mVXhjcxJypq9nfWzftM6fX1pwFsR5lzPmYJu8bmMxGjg_fIj_FBAuSX6Y5RpO8fEyr3QY8jZBDjo4imXKSRP68s8QKnUEVZ330Q1BrbAfLXHjE7Qov2tk3NMonGGrX3eQOkCEN7qjkrBFJ2jI_0Ov1XLnRxm5rDegC41pO6MpkOhoiYfhLyU8QE6465eBWjIsL8GtWhHepPDrjoqPVacwDvBqTqTIaFaFzRZzGj7Zmu4oxyHPxwxq7NpMDo-bv-OYgRNj8kTdCMYgbJQbaiKo359ku1wNj_RbMfJEWq5B67zGNcYH0x44ct7YhRjtuTDfsKUPhjum7SiXkNRPwYemEpVo5vba0JDtLruEZKnrmkiHMjfMU1rvIYk0WqpjmykNdgsWOR0_vHIVIrurFI6wMkFRfbccR-zuB3VNMXR2S5je4wYFxrjqbWMD55pW_wVnTqb5f2a0_tN5vs_wTZgXDzJAVSzOG4VsJiqRgBWRAcPbz-VoX1tftQpg1db3oMd2c_CvKCQpTix0PX84oBKPHB-hlwo0zTfY1K9c-TUwUDjnfuqjOpoFKx2_7itwp8dyZqkNMzw4zYdbZ8HqUdZMryq98EHfq9VChZCyzSrPOJZZczUaDXxWWKZNdpHYLgbm8Tg3gW2W6NYegYC8BfS5EnWzGWL_fpYk9lMkwwq-3p5dfIhr1ZNwgOZXPM5X7CZHxzkxoErIdxzXCzy8h4c4N60HUbuSUFUqGEcA-7mmufvPXm9358byqoMWinYa1mMU68CvFm_-0NCMIhSutJ83fi3FOlQSIoMYHvMmECgNGKzWb0Iz6lHlvMaaGT3OcIjSYEOujSnNd58X6kQAi1Je-llMDYQdQ2sVu2Vbl04z_2ibrE-14qm_vWFwC240VxlTykNmcSShCk_mj35ekMM1mvKYA_6YxYihopYy94M0V9C5DzURit4BN2iaLiNLO5cV57Vijgr6LD460PeQlFKREH6syaQsZnaZAfHq_-CoQy0UNrCM0PVyztAWkEjzQtBkUUIBAJW1-jSH46jKtNmwPAys0HtNJrqRyXRjGg3fGEtOwW5zKnqQfXkI9rCENCp3kpv_8Acbk8_Ari149v.O8cPeH0R7AePTPg1nA4eeQ; _T_ANO=OnfgT0wFpBCQdMmTq6YG0tIOl0Rv2gDsA6nwVYPdZfINHdXcCQteRH8A42Bh8GkGlxjbK4TwgS8of+VJF4nfjHsjj7acxuhG7hgHR/g9SWJ40OwbdgE5bM925AlpE90WrapTVA3uWNr5pZ6ReJSTP44a4QNKl/WfgYMtENf4gwNlBMXFnQnYsB41EbO8grq2deUhhzNHv80MdYthwX77Ebu3X7uRWE4p3ZAIony6Urbml2GafoMt84ZnBbe2e+Gk5+EsbkVUuApPY5U1pjGuE/12Ftv6UpMbbj9L/UmEvDxuqBo4R1VZtSSZfyRDqJETp9A8HKbuiOeclTVZ0po8jA=="
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 64716509

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
